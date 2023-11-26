package main

import (
	"app/config"
	"app/server"
	"flag"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	appConfig = flag.String("config", "config/app.yaml", "application config path")
)

func main() {
	flag.Parse()
	conf, err := config.Parse(*appConfig)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	s, err := server.New(conf, logger)
	if err != nil {
		log.Fatalf("Init server failed: %v", err)
	}
	if err := s.Run(); err != nil {
		logger.Fatalf("Run server failed: %v", err)
	}
}
