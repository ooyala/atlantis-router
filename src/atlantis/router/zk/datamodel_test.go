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

package zk

import (
	"atlantis/router/config"
	"testing"
)

func TestToZkPoolToPool(t *testing.T) {
	conf := config.PoolConfig{
		HealthzEvery:   "1m",
		HealthzTimeout: "9s",
		RequestTimeout: "1s",
		Status:         "OK",
	}

	host1 := config.Host{
		Address: "localhost:8081",
	}
	host2 := config.Host{
		Address: "localhost:8082",
	}
	hosts := map[string]config.Host{
		"host1": host1,
		"host2": host2,
	}

	pool := config.Pool{
		Name:     "test",
		Internal: false,
		Hosts:    hosts,
		Config:   conf,
	}

	zkPool, hosts := ToZkPool(pool)
	recon := zkPool.Pool(hosts)

	if zkPool.Name != "test" || recon.Name != "test" {
		t.Errorf("should preserve name")
	}

	if zkPool.Internal != false || recon.Internal != false {
		t.Errorf("should preserve internal")
	}

	if hosts["host1"].Address != "localhost:8081" || recon.Hosts["host2"].Address != "localhost:8082" {
		t.Errorf("should transform hosts to and fro")
	}

	if zkPool.Config.HealthzEvery != "1m" || recon.Config.RequestTimeout != "1s" {
		t.Errorf("should preserve config")
	}

}

func TestToZkPoolToPoolHttpHeaders(t *testing.T) {
	conf := config.PoolConfig{
		HealthzEvery:   "1m",
		HealthzTimeout: "9s",
		RequestTimeout: "1s",
		Status:         "OK",
	}

	host1 := config.Host{
		Address: "localhost:8081",
	}
	host2 := config.Host{
		Address: "localhost:8082",
	}
	hosts := map[string]config.Host{
		"host1": host1,
		"host2": host2,
	}

	headers := make([]config.HttpHeader, 1)
	headers[0] = config.HttpHeader{Key: "Cache-Control", Value: "max-age:1200"}

	pool := config.Pool{
		Name:     "test",
		Internal: false,
		Hosts:    hosts,
		Config:   conf,
		Headers:  headers,
	}

	zkPool, hosts := ToZkPool(pool)
	recon := zkPool.Pool(hosts)

	if zkPool.Name != "test" || recon.Name != "test" {
		t.Errorf("should preserve name")
	}

	if zkPool.Internal != false || recon.Internal != false {
		t.Errorf("should preserve internal")
	}

	if hosts["host1"].Address != "localhost:8081" || recon.Hosts["host2"].Address != "localhost:8082" {
		t.Errorf("should transform hosts to and fro")
	}

	if zkPool.Config.HealthzEvery != "1m" || recon.Config.RequestTimeout != "1s" {
		t.Errorf("should preserve config")
	}

	if zkPool.Headers[0].Key != "Cache-Control" || zkPool.Headers[0].Value != "max-age:1200" {
		t.Errorf("should preserve headers")
	}

}
