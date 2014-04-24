/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
