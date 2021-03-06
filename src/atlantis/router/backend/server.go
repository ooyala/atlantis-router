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
	"atlantis/router/logger"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Address   string
	Status    ServerStatus
	Metrics   ServerMetrics
	Transport *http.Transport
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
		Status:  NewServerStatus(),
		Metrics: NewServerMetrics(),
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 32,
		},
	}
}

type ResponseError struct {
	Response *http.Response
	Error    error
}

func (s *Server) RoundTrip(req *http.Request, ch chan ResponseError) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			ch <- ResponseError{nil, errors.New(err)}
		}
	}()

	req.URL.Scheme = "http"
	req.URL.Host = s.Address

	resp, err := s.Transport.RoundTrip(req)
	if err == nil {
		ch <- ResponseError{resp, nil}
	} else {
		ch <- ResponseError{resp, err}
	}
}

func (s *Server) Handle(logRecord *logger.HAProxyLogRecord, tout time.Duration) {
	sTime := time.Now()
	s.Metrics.RequestStart()
	defer s.Metrics.RequestDone()

	// X-Forwarded-For; we are a proxy.
	ip := strings.Split(logRecord.Request.RemoteAddr, ":")[0]
	logRecord.Request.Header.Add("X-Forwarded-For", ip)
	logRecord.ServerUpdateRecord(s.Address, s.Metrics.RequestsServiced, s.Metrics.Cost(), sTime)
	resErrCh := make(chan ResponseError)
	tstart := time.Now()
	go s.RoundTrip(logRecord.Request, resErrCh)
	tend := time.Now()
	logRecord.UpdateTr(tstart, tend)
	for {
		select {
		case resErr := <-resErrCh:
			if resErr.Response != nil {
				defer resErr.Response.Body.Close()
			}
			if resErr.Error == nil {
				logRecord.CopyHeaders(resErr.Response.Header)
				logRecord.WriteHeader(resErr.Response.StatusCode)

				err := logRecord.Copy(resErr.Response.Body)
				if err != nil {
					logger.Errorf("[server %s] failed attempting to copy response body: %s\n", s.Address, err)
				} else {
					logRecord.Log()
				}
			} else {
				logger.Errorf("[server %s] failed attempting the roundtrip: %s\n", s.Address, resErr.Error)
				msg := logger.BadGatewayMsg
				status := http.StatusBadGateway
				if strings.Contains(resErr.Error.Error(), "request canceled") {
					msg = logger.GatewayTimeoutMsg
					status = http.StatusGatewayTimeout
				}
				logRecord.Error(msg, status)
				logRecord.Terminate("Server: " + msg)
			}
			return
		case <-time.After(tout):
			// close socket, RoundTrip will return error (or data if the transaction completed before close)
			s.Transport.CancelRequest(logRecord.Request)
		}
	}
}
func (s *Server) CheckStatus(tout time.Duration) {
	r, _ := http.NewRequest("GET", "http://"+s.Address+"/healthz", nil)

	resErrCh := make(chan ResponseError)
	go s.RoundTrip(r, resErrCh)

	for {
		select {
		case resErr := <-resErrCh:
			if resErr.Response != nil {
				defer resErr.Response.Body.Close()
			}
			if resErr.Error == nil {

				//if status has changed then log
				if s.Status.ParseAndSet(resErr.Response) {
					logger.Printf("[server %s] status code changed to %d\n", s.Address, resErr.Response.StatusCode)
				}
			} else {
				//if status has changed then log
				if s.Status.Set(StatusCritical) {
					logger.Errorf("[server %s] status set to critical! : %s\n", s.Address, resErr.Error)
				}
			}
			return
		case <-time.After(tout):
			// close socket, RoundTrip will return error (or data if the transaction completed before close)
			s.Transport.CancelRequest(r)
		}
	}
}

func (s *Server) Cost(accept string) uint32 {
	return s.Status.Cost(accept) + s.Metrics.Cost()
}
