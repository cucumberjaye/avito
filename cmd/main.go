package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("error occured while listen and serve %s", err.Error())
	}
}
