package database

import (
	"errors"
	"fiber-url-shortner/config"
	"fiber-url-shortner/helpers"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database instance
type Dbinstance struct {
	MDb *gorm.DB
	SDb *gorm.DB
}

var DB Dbinstance

var logquery = `
INSERT INTO url_shortner_api_logs (
	user_id, request_id, deviceinfo, apiendpoint, ip_address, geolocation,
	httpmethod, referrer, responsesize, responsetime, statuscode,
	statusmessage, useragent, created_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)`

var insertquery = `INSERT INTO url_shortner (id, url ,created_at) VALUES (?,?,?) on conflict(url)
	do update set url=? returning id`
var selectQuery = `SELECT url FROM url_shortner WHERE id = ? limit 1`
var selectIDQuery = `SELECT count(1) FROM url_shortner WHERE id = ?`

// Connect function
func DBConnect() {
	mdb, err1 := gorm.Open(postgres.Open(config.EnvDBURI("DATABASE_MASTER_URI")), &gorm.Config{})
	sdb, err2 := gorm.Open(postgres.Open(config.EnvDBURI("DATABASE_SLAVE_URI")), &gorm.Config{})
	if err1 != nil || err2 != nil {
		fmt.Printf("database connection failed %s %s", err1.Error(), err2.Error())
		os.Exit(2)
	}
	DB = Dbinstance{
		MDb: mdb,
		SDb: sdb,
	}
}

func DBClose() {
	if dbConnection, err := DB.MDb.DB(); err != nil {
		if err := dbConnection.Ping(); err != nil {
			fmt.Println("ping MDB database failed: " + err.Error())
		}
	} else {
		dbConnection.Close()
		fmt.Println("closed MDB database")
	}
	if dbConnection, err := DB.SDb.DB(); err != nil {
		if err := dbConnection.Ping(); err != nil {
			fmt.Println("ping SDB database failed: " + err.Error())
		}
	} else {
		dbConnection.Close()
		fmt.Println("closed SDB database")
	}
}

func DBPing() (string, error) {
	if dbConnection, err2 := DB.MDb.DB(); err2 != nil {
		return "MDB error", err2
	} else if err := dbConnection.Ping(); err != nil {
		return "MDB ping error", err
	}
	if dbConnection, err2 := DB.SDb.DB(); err2 != nil {
		return "SDB error", err2
	} else if err := dbConnection.Ping(); err != nil {
		return "SDB ping error", err
	}
	return "SDB,MDB: OK", nil
}

func InsertData(id string, url string) (string, error) {
	tx := DB.MDb.Raw(insertquery, id, url, time.Now(), url).Scan(&id)
	return id, tx.Error
}

func CheckIDExists(id string) bool {
	cnt := 0
	if tx := DB.SDb.Raw(selectIDQuery, id).Scan(&cnt); tx.Error != nil {
		log.Fatalf("error: id already exists in db %s, error %s", id, tx.Error.Error())
	}
	return cnt != 0
}

func GetURL(id string) (string, error) {
	url := "url not found"
	tx := DB.SDb.Raw(selectQuery, id).Scan(&url)
	if tx.Error != nil && tx.RowsAffected == 0 {
		return url, errors.New("No data found matching " + id)
	}
	return url, tx.Error
}

func InsertLogData(logs helpers.LogModel) error {
	err := DB.MDb.Exec(logquery, logs.Userid, logs.Requestid, logs.DeviceInfo, logs.Endpoint, logs.IP,
		logs.Location, logs.Method, logs.Referrer, logs.ResponseSize, logs.ResponseTime,
		logs.StatusCode, logs.StatusMessage, logs.UserAgent, logs.CreatedAt).Error
	return err
}
