package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case "GET":
	// 	fmt.Fprintf(w, "Hi, %s!", r.URL.Path[1:])
	case "POST":
		var rec Record

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		rec.Time = time.Now().Format("2006-01-02 15:04:05.000000")

		// db, err := gorm.Open(sqlite.Open("logging_database.db"), &gorm.Config{})
		// db.Table("log_table").Create(&Record{Time: rec.Time, Topic: rec.Topic, Message: rec.Message})

		// if err != nil {
		// 	panic("failed to connect database")
		// }

		// fmt.Println(rec.Time + " " + rec.Topic + " " + rec.Message)
		dbLogChannel <- rec
		localLogChannel <- rec

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are")
	}
}
