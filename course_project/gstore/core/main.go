package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	gServer := NewGServer()
	err := http.ListenAndServe(":4396", gServer)
	if err != nil {
		log.Fatal("cannot listen 4396")
	}
}
