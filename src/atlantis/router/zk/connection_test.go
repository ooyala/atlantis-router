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
	"atlantis/router/testutils"
	"testing"
	"time"
)

func TestManagedZkConn(t *testing.T) {
	// NOTE This test will take ~20s. If that is unacceptably long
	// for your development purposes, uncomment the next 4 lines.
	//
	// if !testing.Verbose() {
	// 	t.Skipf("skipping connection test, use verbose to run")
	// }

	server, err := testutils.NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server")
	}
	defer server.Destroy()

	addr, _ := server.Addr()
	conn := ManagedZkConn(addr)
	defer conn.Shutdown()

	<-conn.ResetCh
	server.Stop()
	server.Start()

	select {
	case <-conn.ResetCh:
		t.Logf("recieved reset on channel")
	case <-time.After(30 * time.Second):
		t.Errorf("connection did not reset")
	}
}
