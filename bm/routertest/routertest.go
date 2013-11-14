package main

import (
	"atlantis/router/backend"
	"log"
	"net/http"
	"time"
)

import _ "net/http/pprof"

var servers = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
}

func main() {
	config := backend.PoolConfig{
		HealthzEvery:   1 * time.Second,
		HealthzTimeout: 1 * time.Second,
		RequestTimeout: 5 * time.Second,
		Status:         "OK",
	}

	/*
		prof, err := os.Create("profile")
		pprof.StartCPUProfile()

		sigINT := make(chan os.Signal, 1)
		signal.Notify(sigINT, os.Interrupt)
		go func(){
			for s := range sigINT{
				pprof.StopCPUProfile()
				os.Exit(0)
			}
		}()
	*/

	pool := backend.NewPool("routertest", config)

	for _, server := range servers {
		pool.AddServer(server, backend.NewServer(server))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pool.Handle(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
