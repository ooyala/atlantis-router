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
	router.New(servers, 8080).Run()
}
