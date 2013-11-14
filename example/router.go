package main

import (
	"atlantis/router/lb"
	"flag"
)

var servers string

func main() {
	flag.StringVar(&servers, "zk", "localhost:2181", "zookeeper connection string")
	balancer := lb.New(servers)
	balancer.Run()
}
