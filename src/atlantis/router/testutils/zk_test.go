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
	"github.com/scalingdata/gozk"
	"testing"
	"time"
)

func TestNewZkServer(t *testing.T) {
	server, err := NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server: %s", err)
	}
	defer server.Destroy()

	addr, _ := server.Addr()
	_, _, err = zookeeper.Dial(addr, 1*time.Second)
	if err != nil {
		t.Errorf("cannot connect to server: %s")
	}
}

func TestNewZkConn(t *testing.T) {
	server, err := NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server: %s", err)
	}
	defer server.Destroy()

	server.Stop()
	_, err = NewZkConn(server, false)
	if err == nil {
		t.Errorf("should fail when server is not running")
	}

	server.Start()
	_, err = NewZkConn(server, false)
	if err != nil {
		t.Errorf("should succeed when server is running")
	}
}

func TestNewZkConnPanic(t *testing.T) {
	server, err := NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server: %s", err)
	}
	defer server.Destroy()

	conn, err := NewZkConn(server, false)
	if err != nil {
		t.Fatalf("cannot connect to zk server")
	}

	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("should panic when connection resets")
			}
		}()
		conn.PanicOnReset()
	}()

	server.Stop()
	time.Sleep(1 * time.Second)

	// It is painful to test NewZkConn(server, true) since the panicing
	// go routine is not in context when it is launched by NewZkConn().
}
