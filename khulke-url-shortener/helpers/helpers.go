package helpers

import (
	"errors"
	"fiber-url-shortner/config"
	"math"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	largestuint64 = 18446744073709551615
)

var whiteListedIPs = strings.Split(config.EnvDBURI("WHITELISTED_IPS"), ",")

func CheckIPAdress(c *fiber.Ctx) error {
	// Add the IPs you want to block here
	path := c.Path()
	if path == "/v1/shorten" {
		var found = false
		clientIP := c.IP()
		for _, whiteListedIP := range whiteListedIPs {
			_, subnet, _ := net.ParseCIDR(whiteListedIP)
			ip := net.ParseIP(clientIP)
			if subnet.Contains(ip){
				found = true
				break
			}
		}
		if !found {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "IP Access Denied for " + clientIP})
		}
	}
	return c.Next()
}

// Allow access to all other IPs

func Base62Encoder(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(5)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}
	return encodedBuilder.String()
}

func Base62Decode(encodedString string) (uint64, error) {
	var number uint64
	length := len(alphabet)

	for i, symbol := range encodedString {
		alphabeticPosition := strings.IndexRune(alphabet, symbol)
		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("cannot find symbol in alphabet")
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func DomainError(url string) bool {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Trim(newURL, " ") // remove all leading/trailing spaces
	newURL = strings.Split(newURL, "/")[0]
	return len(newURL) < 5 || newURL == config.EnvDBURI("URL_SHORTENER_DOMAIN")
}
