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

package backend

import (
	"atlantis/router/testutils"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func newTestConfig() PoolConfig {
	return PoolConfig{
		HealthzEvery:   30 * time.Millisecond,
		HealthzTimeout: 20 * time.Millisecond,
		RequestTimeout: 30 * time.Millisecond,
		Status:         "OK",
	}
}

func TestDummyPool(t *testing.T) {
	pool := DummyPool("dummy")
	defer pool.Shutdown()

	if pool.Dummy != true {
		t.Errorf("should set pool as dummy")
	}
}

func TestNewPool(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	if pool == nil {
		t.Errorf("should create pool")
	}
	if pool.Name != "test" {
		t.Errorf("should set pool name")
	}
	if pool.Dummy != false {
		t.Errorf("should not set dummy")
	}
	if pool.Servers == nil {
		t.Errorf("should init servers")
	}
	if pool.Config != newTestConfig() {
		t.Errorf("should set pool config")
	}
}

func TestAddServer(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	pool.AddServer("127.0.0.1:80", NewServer("127.0.0.1:80"))
	if pool.Servers["127.0.0.1:80"].Address != NewServer("127.0.0.1:80").Address {
		t.Errorf("should add server to pool")
	}

	pool.AddServer("127.0.0.1:80", NewServer("127.0.0.1:80"))
	if len(pool.Servers) > 1 {
		t.Errorf("should de-dup servers in pool")
	}
}

func TestDelServer(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	pool.AddServer("127.0.0.1:80", NewServer("127.0.0.1:80"))
	if len(pool.Servers) != 1 {
		t.Skipf("add server is not working")
	}

	pool.DelServer("127.1.1.1:80")
	if len(pool.Servers) != 1 {
		t.Errorf("should ignore non-existent servers")
	}

	pool.DelServer("127.0.0.1:80")
	if len(pool.Servers) != 0 {
		t.Errorf("should delete server from pool")
	}
}

func TestReconfigure(t *testing.T) {
	pool := NewPool("test", PoolConfig{})
	defer pool.Shutdown()

	pool.Reconfigure(newTestConfig())
	if pool.Config != newTestConfig() {
		t.Errorf("should reconfigure pool config")
	}
}

func TestRunChecks(t *testing.T) {
	conf := newTestConfig()

	pool := NewPool("test", conf)
	defer pool.Shutdown()

	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	backend.SetStatus(http.StatusOK, "MAINTENANCE")
	pool.AddServer(backend.Address(), NewServer(backend.Address()))

	time.Sleep(50 * time.Millisecond)
	if pool.Servers[backend.Address()].Status.Current != StatusMaintenance {
		t.Errorf("should poll for server health")
	}

	conf = PoolConfig{
		HealthzEvery:   2 * time.Second,
		HealthzTimeout: 1 * time.Second,
	}
	pool.Reconfigure(conf)
	time.Sleep(50 * time.Millisecond)

	backend.SetStatus(http.StatusOK, "OK")
	if pool.Servers[backend.Address()].Status.Current != StatusMaintenance {
		t.Errorf("should update check interval")
	}
}

func TestNextMaintenance(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	backend0 := testutils.NewBackend(0, false)
	defer backend0.Shutdown()

	backend1 := testutils.NewBackend(0, false)
	defer backend0.Shutdown()

	backend0.SetStatus(http.StatusOK, "MAINTENANCE")
	backend1.SetStatus(http.StatusOK, "MAINTENANCE")

	pool.AddServer(backend0.Address(), NewServer(backend0.Address()))
	pool.AddServer(backend1.Address(), NewServer(backend1.Address()))
	time.Sleep(50 * time.Millisecond)

	if pool.Next() != nil {
		t.Errorf("should never return server under maintenance")
	}
}

func TestNextCost(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	backend0 := testutils.NewBackend(0, false)
	defer backend0.Shutdown()

	backend1 := testutils.NewBackend(0, false)
	defer backend1.Shutdown()

	pool.AddServer(backend0.Address(), NewServer(backend0.Address()))
	pool.AddServer(backend1.Address(), NewServer(backend1.Address()))
	time.Sleep(50 * time.Millisecond)

	for _, server := range pool.Servers {
		if server.Address == backend0.Address() {
			// leaving a request open for this effect
			// will add to test times
			server.Metrics.RequestStart()
		}
	}

	if pool.Next().Address != backend1.Address() {
		t.Errorf("should return server with least cost")
	}
}

func TestHandleDummy(t *testing.T) {
	pool := DummyPool("test")
	defer pool.Shutdown()

	logRecord, rr := testutils.NewTestHAProxyLogRecord("")
	pool.Handle(logRecord)

	if logRecord.GetResponseStatusCode() != http.StatusBadGateway ||
	     rr.Code != http.StatusBadGateway {
		t.Errorf("should return bad gateway for dummy")
		t.Errorf("%s | %s", logRecord.GetResponseStatusCode(), rr.Code)
	}
}
func TestHandleNoNext(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	logRecord, rr := testutils.NewTestHAProxyLogRecord("")
	pool.Handle(logRecord)

	if logRecord.GetResponseStatusCode() != http.StatusServiceUnavailable ||
	     rr.Code != http.StatusServiceUnavailable {
		t.Errorf("should return unavailable with no next")
	}
}

func TestHandle(t *testing.T) {
	pool := NewPool("test", newTestConfig())
	defer pool.Shutdown()

	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	backend.SetResponse(http.StatusOK, "Mickey Mouse!")

	pool.AddServer(backend.Address(), NewServer(backend.Address()))
	time.Sleep(50 * time.Millisecond)

	logRecord, rr := testutils.NewTestHAProxyLogRecord(backend.URL())
	pool.Handle(logRecord)
	
	body, _ := ioutil.ReadAll(rr.Body)
	if logRecord.GetResponseStatusCode() != http.StatusOK || rr.Code != http.StatusOK ||
	     string(body) != "Mickey Mouse!" {
		t.Errorf("should forward requests to backend")
		t.Errorf("%d | %d | %s", logRecord.GetResponseStatusCode, rr.Code, string(body))
	}
}
