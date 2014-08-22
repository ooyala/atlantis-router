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

func TestNewServer(t *testing.T) {
	server := NewServer("127.0.0.1:80")
	if server.Address != "127.0.0.1:80" {
		t.Errorf("should set server address")
	}

	if server.Status.Current != StatusMaintenance {
		t.Errorf("should mark server under maintenance")
	}
}

func TestRoundTripServerOk(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	req, _ := http.NewRequest("GET", backend.URL(), nil)

	resErrCh := make(chan ResponseError)
	go server.RoundTrip(req, resErrCh)

	resErr := <-resErrCh
	if resErr.Error != nil {
		t.Errorf("should report no error")
	}

	body, _ := ioutil.ReadAll(resErr.Response.Body)
	if string(body) != "testutils backend" {
		t.Errorf("should return server's response")
	}

	return
}

func TestRoundTripServerError(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	backend.SetResponse(http.StatusInternalServerError, "No shit, Sherlock!")
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	req, _ := http.NewRequest("GET", backend.URL(), nil)

	resErrCh := make(chan ResponseError)
	go server.RoundTrip(req, resErrCh)

	resErr := <-resErrCh
	if resErr.Error != nil {
		// transport must not care about server's response
		// as long as there are no connection errors
		t.Errorf("should report no error")
	}

	body, _ := ioutil.ReadAll(resErr.Response.Body)
	if string(body) != "No shit, Sherlock!" {
		t.Errorf("should return server's response")
	}

	return
}

func TestHandleResponse(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	backend.SetResponse(http.StatusOK, "The eagle has landed.")
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	logRecord, rr := testutils.NewTestHAProxyLogRecord(backend.URL())
	server.Handle(logRecord, 100*time.Millisecond)

	if logRecord.GetResponseStatusCode() != http.StatusOK {
		t.Errorf("should set status code")
	}
	body, _ := ioutil.ReadAll(rr.Body)
	if string(body) != "The eagle has landed." {
		t.Errorf("should set response body")
	}
}

func TestHandleXForwardedFor(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	logRecord, _ := testutils.NewTestHAProxyLogRecord(backend.URL())
	server.Handle(logRecord, 100*time.Millisecond)

	elm := backend.Handler.Recorded.Front()
	rec := elm.Value.(testutils.RequestAndTime).R
	if rec.Header["X-Forwarded-For"] == nil {
		t.Errorf("should set x-forwarded-for")
	}
}

func TestHandleResponseHeaders(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	logRecord, _ := testutils.NewTestHAProxyLogRecord(backend.URL() + "/healthz")
	server.Handle(logRecord, 100*time.Millisecond)

	if logRecord.GetResponseHeaders()["Server-Status"][0] != "OK" {
		t.Errorf("should copy response headers")
	}
}

func TestHandleTimeout(t *testing.T) {
	backend := testutils.NewBackend(100, false)
	backend.SetResponse(http.StatusOK, "Ba-ba-ba-ba-ba-na-na!")
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	logRecord, rr := testutils.NewTestHAProxyLogRecord(backend.URL())
	server.Handle(logRecord, 10*time.Millisecond)

	if logRecord.GetResponseStatusCode() != http.StatusGatewayTimeout ||
		rr.Code != http.StatusGatewayTimeout {
		t.Errorf("should report status code 502")
	}

	body, _ := ioutil.ReadAll(rr.Body)
	if string(body) != "Gateway Timeout\n" {
		t.Errorf("should report 'Gateway Timeout'")
	}
}

func TestCheckStatus(t *testing.T) {
	backend := testutils.NewBackend(0, false)
	defer backend.Shutdown()

	server := NewServer(backend.Address())

	backend.SetStatus(http.StatusInternalServerError, "UNKNOWN")
	server.CheckStatus(100 * time.Millisecond)
	if server.Status.Current != StatusMaintenance {
		t.Errorf("should set status to maintenance on error")
	}

	backend.SetStatus(http.StatusOK, "DEGRADED")
	server.CheckStatus(100 * time.Millisecond)
	if server.Status.Current != StatusDegraded {
		t.Errorf("should parse response headers on 200s")
	}
}

func TestCheckStatusTimeout(t *testing.T) {
	backend := testutils.NewBackend(100, false)
	defer backend.Shutdown()

	server := NewServer(backend.Address())
	server.CheckStatus(10 * time.Millisecond)
	if server.Status.Current != StatusCritical {
		t.Errorf("should set status to critical on timeout")
	}
}
