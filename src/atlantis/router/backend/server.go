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
	"strconv"
	"strings"
	"time"
	"os"
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

func (s *Server) logPrefix(r *http.Request, tstart time.Time) string {
	now := time.Now()

	var rtt0, rtt1 int64

	// Calculate the total round trip time if front end inserted the atlantis-arrival-time
	// header before routing the request. The header is assumed to be created by calling
	// time.Now().UnixNano() or equivalent.
	arr, err := strconv.ParseInt(r.Header.Get("atlantis-arrival-time"), 10, 64)
	if err != nil {
		rtt0 = now.UnixNano() - arr
	}

	rtt1 = now.UnixNano() - tstart.UnixNano()

	// Log prefix includes server address, request source and URI, and round trip times.
	return fmt.Sprintf("[server %s][request %s|%s][rtt %d|%d]", s.Address, r.RemoteAddr,
		r.URL, rtt0, rtt1)
}

func (s *Server) Handle(logRecord *logger.HAProxyLogRecord, tout time.Duration) {
	s.Metrics.RequestStart()
	defer s.Metrics.RequestDone()

	// X-Forwarded-For; we are a proxy.
	ip := strings.Split(logRecord.Request.RemoteAddr, ":")[0]
	logRecord.Request.Header.Add("X-Forwarded-For", ip)
	logRecord.ServerUpdateRecord("", s.Metrics.RequestsServiced, s.Metrics.Cost())
	resErrCh := make(chan ResponseError)
	tstart := time.Now()
	go s.RoundTrip(logRecord.Request, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			logger.Printf("%s %d", s.logPrefix(logRecord.Request, tstart), resErr.Response.StatusCode)
			defer resErr.Response.Body.Close()

			logRecord.CopyHeaders(resErr.Response.Header)
			logRecord.SetResponseStatusCode(resErr.Response.StatusCode)	

			err := logRecord.Copy(resErr.Response.Body)
			if err != nil {
				logger.Errorf("%s %s", s.logPrefix(logRecord.Request, tstart), err)
			} else {
				logRecord.Log(os.Stdout)
			}
		} else {
			logger.Errorf("%s %s", s.logPrefix(logRecord.Request, tstart), resErr.Error)
			logRecord.Error(logger.BadGatewayMsg, http.StatusBadGateway)
		}
	case <-time.After(tout):
		s.Transport.CancelRequest(logRecord.Request)
		logger.Printf("%s timeout", s.logPrefix(logRecord.Request, tstart))
		logRecord.Error(logger.GatewayTimeoutMsg, http.StatusGatewayTimeout)
	}
}
func (s *Server) CheckStatus(tout time.Duration) {
	r, _ := http.NewRequest("GET", "http://"+s.Address+"/healthz", nil)

	resErrCh := make(chan ResponseError)
	tstart := time.Now()
	go s.RoundTrip(r, resErrCh)

	select {
	case resErr := <-resErrCh:
		if resErr.Error == nil {
			logger.Printf("%s %d", s.logPrefix(r, tstart), resErr.Response.StatusCode)
			defer resErr.Response.Body.Close()
			s.Status.ParseAndSet(resErr.Response)
		} else {
			logger.Errorf("%s %s", s.logPrefix(r, tstart), resErr.Error)
			s.Status.Set(StatusCritical)
		}
	case <-time.After(tout):
		s.Transport.CancelRequest(r)
		logger.Errorf("%s timeout", s.logPrefix(r, tstart))
		s.Status.Set(StatusCritical)
	}
}

func (s *Server) Cost(accept string) uint32 {
	return s.Status.Cost(accept) + s.Metrics.Cost()
}

func (s *Server) Shutdown() {
//TODO previously they shut down the copier here
//but now that the copier is a statics to the HAProxyLog package 
//dunno where exactly to shut down probably when router shuts down
}
