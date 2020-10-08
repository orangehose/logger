package main

import (
	"fmt"
	"log"
	"net/http"
)

var localLogChannel chan Record
var dbLogChannel chan Record

func init() {
	localLogChannel = make(chan Record, 100)
	dbLogChannel = make(chan Record, 100)

	go writeToDbLogRoutine(dbLogChannel)
	go writeToLocalLogRoutine(localLogChannel)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/log", handler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
