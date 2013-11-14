package testutils

import (
	"errors"
	"launchpad.net/gozk"
	"os"
	"time"
)

const (
	macOsXPath = "/usr/local/Cellar/zookeeper/3.4.5/libexec"
	linuxPath  = "/usr/share/java"
	serverPort = 2182
	serverDir  = "/tmp/zkserver"
)

func NewZkServer() (*zookeeper.Server, error) {
	os.RemoveAll(serverDir)

	path := macOsXPath
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		path = linuxPath
	}

	server, err := zookeeper.CreateServer(serverPort, serverDir, path)
	if err != nil {
		return nil, err
	}
	err = server.Start()
	if err != nil {
		return nil, err
	}

	return server, nil
}

type ZkConn struct {
	Conn *zookeeper.Conn
	evCh <-chan zookeeper.Event
}

func NewZkConn(server *zookeeper.Server, panicOnReset bool) (*ZkConn, error) {
	addr, _ := server.Addr()
	conn, eventCh, err := zookeeper.Dial(addr, 1*time.Second)
	if err != nil {
		return nil, err
	}

	zkConn := &ZkConn{
		Conn: conn,
		evCh: eventCh,
	}

	tout := time.After(5 * time.Second)
	for {
		select {
		case ev := <-eventCh:
			if ev.State == zookeeper.STATE_CONNECTED {
				if panicOnReset {
					go zkConn.PanicOnReset()
				}
				return zkConn, nil
			}
		case <-tout:
			return nil, errors.New("timeout connecting to zookeeper")
		}
	}
}

func (z *ZkConn) PanicOnReset() {
	for {
		ev := <-z.evCh
		if ev.State == zookeeper.STATE_EXPIRED_SESSION ||
			ev.State == zookeeper.STATE_CONNECTING {
			panic("zookeeper connection lost")
		} else if ev.State == zookeeper.STATE_CLOSED {
			// connection closed by tests
			return
		}
	}
}
