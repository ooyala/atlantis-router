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

package testutils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestStatus(t *testing.T) {
	backend := NewBackend(0, false)
	defer backend.Shutdown()

	req, _ := http.NewRequest("GET", backend.URL()+"/healthz", nil)

	client := &http.Client{}

	res, _ := client.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("should default status code to http.StatusOK")
	}
	if res.Header.Get("Server-Status") != "OK" {
		t.Errorf("should default Server-Status header to OK")
	}

	backend.SetStatus(http.StatusInternalServerError, "UNKNOWN")

	res, _ = client.Do(req)
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("should set status code to http.StatusInternalServerError")
	}
	if res.Header.Get("Server-Status") != "UNKNOWN" {
		t.Errorf("should set Server-Status header to UNKNOWN")
	}

	return
}

func TestResponse(t *testing.T) {
	backend := NewBackend(0, false)
	defer backend.Shutdown()

	req, _ := http.NewRequest("GET", backend.URL(), nil)

	client := &http.Client{}

	res, _ := client.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("should default response code to http.StatusOK")
	}
	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != "testutils backend" {
		t.Errorf("should default response body to 'testutils backend'")
	}

	backend.SetResponse(http.StatusInternalServerError, "Holy Shit Batman!")

	res, _ = client.Do(req)
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("should set response code to http.StatusInternalServerError")
	}
	body, _ = ioutil.ReadAll(res.Body)
	if string(body) != "Holy Shit Batman!" {
		t.Errorf("should set response body to 'Holy Shit Batman!'")
	}
}

func TestAddress(t *testing.T) {
	backend := NewBackend(0, false)
	defer backend.Shutdown()

	url := fmt.Sprintf("http://%s", backend.Address())
	if url != backend.URL() {
		t.Errorf("address and url should point to the backend")
	}
}

func TestDelays(t *testing.T) {
	l := 10
	N := 10

	backend := NewBackend(l, true)
	defer backend.Shutdown()

	req, _ := http.NewRequest("GET", backend.URL(), nil)

	client := &http.Client{}

	start := time.Now()
	for i := 0; i < N; i++ {
		client.Do(req)
	}
	delayNs := float64(time.Now().UnixNano()-start.UnixNano()) / float64(N)
	delayMs := int(delayNs / float64(delayNs))

	sigma := float64(l / N) // is this the correct convergence rate?
	if float64(delayMs-l) > sigma {
		t.Errorf("should average around %d vs %d", l, delayMs)
	}

	return
}

func TestRecording(t *testing.T) {
	backend := NewBackend(0, false)
	defer backend.Shutdown()

	client := &http.Client{}
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("%s/%d", backend.URL(), i)
		req, _ := http.NewRequest("GET", url, nil)
		client.Do(req)
	}

	e := backend.Handler.Recorded.Front()
	for i := 0; i < 10; i++ {
		reqURL := fmt.Sprintf("/%d", i)
		reqRec := e.Value.(RequestAndTime).R.URL.Path
		if reqRec != reqURL {
			t.Errorf("played %s, recorded %s", reqURL, reqRec)
		}
		e = e.Next()
	}

	return
}
