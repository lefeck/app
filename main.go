package main

import (
	"app/config"
	"app/server"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	s, err := server.New(conf)
	if err != nil {
		log.Fatalf("Init server failed: %v", err)
	}
	s.Run()
}
