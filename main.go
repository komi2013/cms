package main

import (
	"fmt"
	"log"
	"net/http"

	"./app"
	"./common"
)

func main() {
	http.HandleFunc("/htm/", app.Htm)
	http.HandleFunc("/", app.Top)
	fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))
}