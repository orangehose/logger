package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database string = "logging_database.db"
var tableName string = "log_table"

func writeToLocalLogRoutine(channel chan Record) {
	var rec Record
	for {
		rec = <-channel
		msg := rec.Time + " " + rec.Topic + " " + rec.Message
		writeToLocalLog(msg)
	}
}

func writeToDbLogRoutine(channel chan Record) {
	var rec Record
	var isDbExists bool

	for {
		rec = <-channel
		isDbExists = fileExists(database)

		if rec != (Record{}) && isDbExists == true {
			db, err := initDbSession()
			if err != nil {
				// trying again
				tryToWrite(rec)
			} else {
				db.Table(tableName).Create(&Record{Time: rec.Time, Topic: rec.Topic, Message: rec.Message})
			}
		} else {
			// trying again
			tryToWrite(rec)
		}
	}
}

func tryToWrite(rec Record) {
	for {
		if fileExists(database) == true {
			db, err := initDbSession()
			if err == nil {
				db.Table(tableName).Create(&Record{Time: rec.Time, Topic: rec.Topic, Message: rec.Message})
				break
			}
		}
	}
}

func initDbSession() (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(database), &gorm.Config{})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
