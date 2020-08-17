package main

import (
	"fixator/fixator"
	"fixator/handler"
	"log"
	"net/http"
)

func main() {
	s := &fixator.Fixator{}
	api := handler.New(s)

	log.Fatal(http.ListenAndServe(":3030", api.Router()))
}
