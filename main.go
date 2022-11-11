package main

import (
	"log"
	"net/http"
	"tour-le-shit-go/pkg/server"
)

const port = ":4000"

func main() {
	srv := server.New()
	log.Fatal(http.ListenAndServe(port, srv))
}
