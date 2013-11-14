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
