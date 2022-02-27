package main

import (
	"fmt"
	"log"
	"net/http"

	"./app"
	"./common"
)

func main() {
	// http.HandleFunc("/htm/", app.Htm)
	http.HandleFunc("/advertisement/", app.Advertisement)
	http.HandleFunc("/category/", app.Category)
	http.HandleFunc("/note/", app.Note)
	fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))
}
