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
	"container/list"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

type Status struct {
	Code   int
	Header string
}

func (s *Status) Set(code int, header string) {
	s.Code = code
	s.Header = header
}

type Response struct {
	Code int
	Body string
}

func (r *Response) Set(code int, body string) {
	r.Code = code
	r.Body = body
}

type RequestAndTime struct {
	R http.Request
	T time.Time
}

type Handler struct {
	MeanWait int
	PowerLaw bool
	Response Response
	Status   Status
	Recorded *list.List
}

func NewHandler(meanWaitMs int, powerLaw bool) *Handler {
	if meanWaitMs == 0 {
		// a reasonable default
		meanWaitMs = 10
	}

	handler := &Handler{
		MeanWait: meanWaitMs,
		PowerLaw: powerLaw,
		Response: Response{http.StatusOK, "testutils backend"},
		Status:   Status{http.StatusOK, "OK"},
		Recorded: list.New(),
	}

	return handler
}

func (h *Handler) Wait() {
	var waitMs int
	if h.PowerLaw {
		waitMs = int(rand.ExpFloat64() * float64(h.MeanWait))
	} else {
		waitMs = h.MeanWait
	}

	time.Sleep(time.Duration(waitMs) * time.Millisecond)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Recorded.PushBack(RequestAndTime{*r, time.Now()})
	h.Wait()
	h.deMux(w, r)
}

func (h *Handler) deMux(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/healthz") {
		w.Header().Set("Server-Status", h.Status.Header)
		w.WriteHeader(h.Status.Code)
		w.Write([]byte("testutils healthz\n"))
	} else {
		w.WriteHeader(h.Response.Code)
		w.Write([]byte(h.Response.Body))
	}
}

type Backend struct {
	Server  *httptest.Server
	Handler *Handler
}

func NewBackend(meanWaitMs int, powerLaw bool) *Backend {
	backend := &Backend{}
	backend.Handler = NewHandler(meanWaitMs, powerLaw)
	backend.Server = httptest.NewServer(backend.Handler)
	return backend
}

func (b *Backend) SetResponse(code int, body string) {
	b.Handler.Response.Set(code, body)
}

func (b *Backend) SetStatus(code int, header string) {
	b.Handler.Status.Set(code, strings.ToUpper(header))
}

func (b *Backend) URL() string {
	return b.Server.URL
}

func (b *Backend) Address() string {
	return b.URL()[7:]
}

func (b *Backend) Shutdown() {
	b.Server.Close()
}
