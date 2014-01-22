package main

import (
	"atlantis/router/router"
	"flag"
	"log"
	"log/syslog"
)

var servers string

func main() {
	// Logging to syslog is more performant, which matters.
	w, err := syslog.New(syslog.LOG_INFO, "atlantis-router")
	if err != nil {
		log.Println("[ERROR] cannot log to syslog!")
	} else {
		log.SetOutput(w)
		log.SetFlags(0)
	}

	flag.StringVar(&servers, "zk", "localhost:2181", "zookeeper connection string")
	router.New(servers).Run()
}
