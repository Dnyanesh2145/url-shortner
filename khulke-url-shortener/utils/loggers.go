package utils

import (
	"fiber-url-shortner/database"
	"fiber-url-shortner/helpers"
	"fmt"
	"io"
	"net"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var placeholders = []string{
	"Timestamp", "IP", "requestid", "userId", "Method",
	"Endpoint", "StatusCode", "ResponseSize", "ResponseTime", "UserAgent",
	"DeviceInfo", "Location", "referrer", "StatusMessage",
}

const (
	CONFIG_SERVICENAME = "shortener_service"
	CONFIG_LOGFORMAT   = "${time} <----> ${ip} <----> ${locals:reqid} <----> ${locals:userId} <----> ${method}" +
		" <----> ${path} <----> ${status} <----> ${bytesSent} <----> ${latency} <----> ${ua} <----> ${locals:device}" +
		" <----> ${locals:latlong} <----> ${referrer} <----> ${locals:msg}"
	CONFIG_LOGTIME_FORMAT   = "2006-01-02T15:04:05.000"
	CONFIG_DEFAULT_REQID    = "000000"
	CONFIG_DEFAULT_UID      = "ffffff"
	CONFIG_TIMESTAMP_FORMAT = "2006-01-02T15:04:05.999"
)

var DEBUG_LOGFORMAT = regexp.MustCompile(`(?U:[^( <\----> )]+?)`).ReplaceAllString(CONFIG_LOGFORMAT, "%s")

type LogStringChannelWriter chan string
type DebugLogStringChannelWriter chan string

var LogStringChannel LogStringChannelWriter
var DebugLogStringChannel DebugLogStringChannelWriter

func (ch LogStringChannelWriter) Write(p []byte) (l int, err error) {
	l = len(p)
	if l < 1 {
		return l, io.ErrShortWrite
	}
	ch <- string(p)
	return l, nil
}

func (ch DebugLogStringChannelWriter) Write(p []byte) (l int, err error) {
	l = len(p)
	if l < 1 {
		return l, io.ErrShortWrite
	}
	ch <- string(p)
	return l, nil
}

func CloseLogStringChannel() {
	close(LogStringChannel)
	fmt.Println("closed logstringchannel")
}

func GetLocalIP() string {
	if conn, err := net.Dial("ip:icmp", "google.com"); err != nil {
		fmt.Println("unable to get localip", err.Error())
		return "::1"
	} else {
		return conn.LocalAddr().String()

	}
}

// Puts logs to channel then to redis, other service will dequeue from redis to logsdb
func DebugAPI(c *fiber.Ctx, msg string) {
	pc, file, no, _ := runtime.Caller(1) // pc, file, no, ok
	debugLogString := fmt.Sprintf(DEBUG_LOGFORMAT,
		time.Now().Format(CONFIG_LOGTIME_FORMAT), c.IP(), c.Context().Value("reqid"),
		c.Context().Value("uid"), "DEBUG", runtime.FuncForPC(pc).Name(), file, strconv.Itoa(no), "0",
		"0", c.Request().Header.UserAgent(), c.Context().Value("device"), "DEV", msg)
	DebugLogStringChannel <- debugLogString
	fmt.Print(debugLogString)
}

// Puts logs to channel then to redis, other service will dequeue from redis to logsdb
func DebugCMD(latency, msg string) {
	pc, file, no, _ := runtime.Caller(1) // pc, file, no, ok
	debugLogString := fmt.Sprintf("%s <----> %s <----> %s <----> %s <----> %s <----> %s <----> %d <----> %s<---->%s<---->%s<---->%s<---->%s<---->%s<---->%s",
		time.Now().Format(CONFIG_LOGTIME_FORMAT), GetLocalIP(), "", "", runtime.FuncForPC(pc).Name(),
		file, no, "", latency, "", "", "", "", msg)
	DebugLogStringChannel <- debugLogString
	fmt.Println(debugLogString)
}

func LogStringChannelConsumer() {
	fmt.Println("started LogStringChannelConsumer")
	for logMsg := range LogStringChannel {
		result := make(map[string]string)
		// Split the log entry into parts
		parts := strings.Split(logMsg, " <----> ")

		// Iterate through placeholders and corresponding parts
		for i, placeholder := range placeholders {
			if i < len(parts) {
				result[placeholder] = parts[i]
			}
		}

		statuscode, err := strconv.Atoi(result["StatusCode"])

		if err != nil {
			fmt.Println("str not conveeted to int64", err.Error())
		}

		responsesize, err := strconv.Atoi(result["ResponseSize"])

		if err != nil {
			fmt.Println("str not conveeted to int", err.Error())
		}

		h, _ := time.ParseDuration(result["ResponseTime"])

		// Use the appropriate layout for your timestamp format
		timestamp, err := time.Parse(CONFIG_TIMESTAMP_FORMAT, result["Timestamp"])
		if err != nil {
			fmt.Println("str not conveeted to timestamp", err.Error())
		}

		logs := helpers.LogModel{
			Userid:        result["userId"],
			Requestid:     result["requestid"],
			DeviceInfo:    result["DeviceInfo"],
			Endpoint:      result["Endpoint"],
			IP:            result["IP"],
			Location:      result["Location"],
			Method:        result["Method"],
			Referrer:      result["Referrer"],
			ResponseSize:  responsesize,
			ResponseTime:  int(h.Milliseconds()),
			StatusCode:    statuscode,
			StatusMessage: result["StatusMessage"],
			UserAgent:     result["UserAgent"],
			CreatedAt:     timestamp,
		}

		res := database.InsertLogData(logs)
		if err != nil {
			fmt.Println("error", res)
		}
	}
}

func DebugLogStringChannelConsumer() {
	fmt.Println("started DebugLogStringChannelConsumer")
	key := "dequeue_insert_table_log_" + CONFIG_SERVICENAME
	for logMsg := range DebugLogStringChannel {

		fmt.Println(logMsg, key)
	}
}
