package main

import (
	"log"
	"net/http"

	_ "github.com/bearprada/who_most_like_you"
)

func main() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("%v", err)
	}
}
