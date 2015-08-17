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
	"atlantis/router/logger"
	"errors"
	"github.com/scalingdata/gozk"
	"sync"
	"time"
)

type ZkConn struct {
	sync.Mutex
	ResetCh   chan bool
	servers   string
	Conn      *zookeeper.Conn
	eventCh   <-chan zookeeper.Event
	killCh    chan bool
	connected bool
}

func ManagedZkConn(servers string) *ZkConn {
	zk := &ZkConn{
		ResetCh:   make(chan bool),
		servers:   servers,
		killCh:    make(chan bool),
		connected: false,
	}

	go zk.dialExclusive()

	return zk
}

func (z *ZkConn) Shutdown() {
	z.killCh <- true
	z.Conn.Close()
}

func (z *ZkConn) dialExclusive() {
	z.Lock()

	for err := z.dial(); err != nil; {
		logger.Printf("[zkconn %p] z.dial(): %s", z, err)
	}

	z.Unlock()

	z.ResetCh <- true
}

func (z *ZkConn) dial() error {
	var err error
	z.connected = false
	z.Conn, z.eventCh, err = zookeeper.Dial(z.servers, 30*time.Second)
	if err != nil {
		return err
	}

	err = z.waitOnConnect()
	if err != nil {
		return err
	}
	z.connected = true

	go z.monitorEventCh()

	return nil
}

func (z *ZkConn) waitOnConnect() error {
	for {
		ev := <-z.eventCh
		logger.Printf("[zkconn %p] eventCh-> %d %s in waitOnConnect", z, ev.State, ev)

		switch ev.State {
		case zookeeper.STATE_CONNECTED:
			return nil
		case zookeeper.STATE_CONNECTING:
			continue
		default:
			return errors.New(ev.String())
		}
	}
}

func (z *ZkConn) monitorEventCh() {
	for {
		select {
		case ev := <-z.eventCh:
			logger.Printf("[zkconn %p] eventCh -> %d %s in monitorEventCh", z, ev.State, ev)
			if ev.State == zookeeper.STATE_EXPIRED_SESSION ||
				ev.State == zookeeper.STATE_CONNECTING {
				z.dialExclusive()
				return
			}

		case <-z.killCh:
			return
		}
	}
}

func (z *ZkConn) IsConnected() bool {
	if z.connected {
		return true
	} else {
		return false
	}
}
