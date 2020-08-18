package main

import (
	"log"
	"net/http"

	"fixator/config"
	"fixator/fixator"
	"fixator/handler"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	s := fixator.New(conf.Fixator)
	api := handler.New(s, conf.Service)

	log.Fatal(http.ListenAndServe(conf.Host+":"+conf.Port, api.Router()))
}
