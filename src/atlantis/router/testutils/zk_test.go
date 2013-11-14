package testutils

import (
	"launchpad.net/gozk"
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
