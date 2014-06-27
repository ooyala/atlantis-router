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
	"sync/atomic"
)

type ServerMetrics struct {
	RequestsInFlight uint32
	RequestsServiced uint64
}

func NewServerMetrics() ServerMetrics {
	return ServerMetrics{
		RequestsInFlight: 0,
		RequestsServiced: 0,
	}
}

func (s *ServerMetrics) RequestStart() {
	atomic.AddUint32(&s.RequestsInFlight, uint32(1))
	s.RequestsServiced++
}

func (s *ServerMetrics) RequestDone() {
	atomic.AddUint32(&s.RequestsInFlight, ^uint32(0))
}

func (s *ServerMetrics) Cost() uint32 {
	return s.RequestsInFlight
}

type ConMetrics struct {
	ActiveConnections uint32
	TotalConnections  uint64
}

func NewConMetrics() ConMetrics {
	return ConMetrics{
		ActiveConnections: 0,
		TotalConnections:  0,
	}
}

func (c *ConMetrics) ConnectionStart() {
	atomic.AddUint32(&c.ActiveConnections, uint32(1))
	c.TotalConnections++
}

func (c *ConMetrics) ConnectionDone() {
	atomic.AddUint32(&c.ActiveConnections, ^uint32(0))
}

func (c *ConMetrics) GetActiveConnections() uint32 {
	return c.ActiveConnections
}

func (c *ConMetrics) GetTotalConnections() uint64 {
	return c.TotalConnections
}
