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
	"github.com/scalingdata/gozk"
	"log"
	"runtime"
	"testing"
	"time"
)

const (
	StatusInvalid = iota
	StatusCreated
	StatusChanged
	StatusDeleted
)

type Track struct {
	Status int
	Change int
}

type TrackingCallbacks struct {
	Track map[string]*Track
}

func (c TrackingCallbacks) Created(path, json string) {
	if val, ok := c.Track[path]; !ok {
		c.Track[path] = &Track{StatusCreated, 0}
	} else if val.Status == StatusDeleted {
		val.Status = StatusCreated
	} else {
		val.Status = StatusInvalid
	}
}

func (c TrackingCallbacks) Changed(path, json string) {
	if val, ok := c.Track[path]; !ok {
		c.Track[path] = &Track{StatusInvalid, 0}
	} else if val.Status == StatusDeleted {
		val.Status = StatusInvalid
	} else {
		val.Status = StatusChanged
		val.Change += 1
	}
}

func (c TrackingCallbacks) Deleted(path string) {
	if val, ok := c.Track[path]; !ok {
		c.Track[path] = &Track{StatusInvalid, 0}
	} else if val.Status == StatusDeleted {
		val.Status = StatusInvalid
	} else {
		val.Status = StatusDeleted
	}
}

func (c TrackingCallbacks) valid() bool {
	for path, node := range c.Track {
		if node.Status == StatusInvalid {
			log.Printf("%s is invalid", path)
			return false
		}
	}
	return true
}

func (c TrackingCallbacks) tally(nodes []string, status int) bool {
	for _, node := range nodes {
		if val, ok := c.Track[node]; !ok {
			return false
		} else if val.Status != status {
			return false
		}
	}
	return true
}

func TestManageTree(t *testing.T) {
	server, err := testutils.NewZkServer()
	if err != nil {
		t.Fatalf("cannot start zk server")
	}
	addr, _ := server.Addr()

	numGoRoutine := runtime.NumGoroutine()

	zk := ManagedZkConn(addr)
	<-zk.ResetCh
	zk.Conn.Create("/testing", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	setupList := []string{
		"/testing/setup0",
		"/testing/setup1",
	}
	for _, node := range setupList {
		zk.Conn.Create(node, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	}

	track := TrackingCallbacks{
		Track: make(map[string]*Track),
	}
	go zk.ManageTree("/testing", track, track, track)
	time.Sleep(100 * time.Millisecond)

	if len(track.Track) != 2 {
		t.Errorf("should be tracking 2 nodes")
	} else if !track.tally(setupList, StatusCreated) {
		t.Errorf("setupList should have status created")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	level0 := []string{
		"/testing/node0",
		"/testing/node1",
		"/testing/node2",
	}
	for _, node := range level0 {
		go zk.Conn.Create(node, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	}
	time.Sleep(100 * time.Millisecond)

	if len(track.Track) != 5 {
		t.Errorf("should be tracking 2 nodes")
	} else if !track.tally(level0, StatusCreated) {
		t.Errorf("level0 should have status created")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	level1 := []string{
		"/testing/node0/node0",
		"/testing/node0/node1",
		"/testing/node0/node2",
	}
	for _, node := range level1 {
		go zk.Conn.Create(node, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	}
	time.Sleep(100 * time.Millisecond)

	if len(track.Track) != 8 {
		t.Errorf("should be tracking 6 nodes")
	} else if !track.tally(level0, StatusCreated) {
		t.Errorf("level0 should have status created")
	} else if !track.tally(level1, StatusCreated) {
		t.Errorf("level1 should have status created")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	level2 := []string{
		"/testing/node0/node0/node0",
		"/testing/node0/node0/node1",
	}
	for _, node := range level2 {
		go zk.Conn.Create(node, "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	}
	time.Sleep(100 * time.Millisecond)

	if len(track.Track) != 10 {
		t.Errorf("should be tracking 8 nodes")
	} else if !track.tally(level0, StatusCreated) {
		t.Errorf("level0 should have status created")
	} else if !track.tally(level1, StatusCreated) {
		t.Errorf("level1 should have status created")
	} else if !track.tally(level2, StatusCreated) {
		t.Errorf("level2 should have status created")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	changeList := []string{
		"/testing/node1",
		"/testing/node0/node0",
		"/testing/node0/node0/node1",
	}
	for _, node := range changeList {
		go zk.Conn.Set(node, "changed!", -1)
	}
	time.Sleep(100 * time.Millisecond)

	if !track.tally(changeList, StatusChanged) {
		t.Errorf("changeList should have status changed")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	deleteList := []string{
		"/testing/node0/node0/node0",
		"/testing/node0/node0/node1",
		"/testing/node0/node0",
	}
	for _, node := range deleteList {
		zk.Conn.Delete(node, -1)
	}

	if !track.tally(deleteList, StatusDeleted) {
		t.Errorf("deleteList should have status deleted")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	changeList2 := []string{
		"/testing/node0/node1",
		"/testing/node2",
	}
	for _, node := range changeList2 {
		go zk.Conn.Set(node, "again!", -1)
	}
	time.Sleep(100 * time.Millisecond)

	if !track.tally(changeList2, StatusChanged) {
		t.Errorf("changeList2 should have status changed")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	zk.Conn.Create("/testing/node0/node0", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	time.Sleep(100 * time.Millisecond)

	if !track.tally([]string{"/testing/node0/node0"}, StatusCreated) {
		t.Errorf("/testing/node0/node0 should have status created")
	}
	if !track.valid() {
		t.Errorf("should be valid")
	}

	zk.Shutdown()

	// This test consistently fails under gocov and passes under
	// go test. Switch false to true if you need to verify that
	// ManageTree is not leaking goroutines.
	if false {
		time.Sleep(100 * time.Millisecond)
		if runtime.NumGoroutine() != numGoRoutine {
			t.Errorf("should not leak go routines")
		}
	}

	server.Destroy()
}
