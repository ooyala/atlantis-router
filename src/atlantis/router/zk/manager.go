package zk

import (
	"launchpad.net/gozk"
	"log"
	"path"
)

type EventCallbacks interface {
	Created(path, json string)
	Deleted(path string)
	Changed(path, json string)
}

func (z *ZkConn) ManageNode(node string, callbacks EventCallbacks) error {
	content, _, eventCh, err := z.Conn.GetW(node)
	if err != nil {
		log.Printf("GetW(%s): %s", node, err)
		return err
	}

	callbacks.Created(node, content)

	go func() {
		for {
			ev := <-eventCh

			if ev.State == zookeeper.STATE_CLOSED {
				// shutdown was called on ZkConn?
				return
			}

			if ev.State == zookeeper.STATE_EXPIRED_SESSION ||
				ev.State == zookeeper.STATE_CONNECTING {
				log.Printf("GetW(%s): exiting on loss of connection", node)
				return
			}

			switch ev.Type {
			case zookeeper.EVENT_DELETED:
				callbacks.Deleted(node)
				return
			case zookeeper.EVENT_CHANGED:
				content, _, eventCh, err = z.Conn.GetW(node)
				if err != nil {
					log.Printf("GetW(%s): %s", node, err)
					return
				}
				callbacks.Changed(node, content)
			}
		}
	}()

	return nil
}

func (z *ZkConn) ManageTree(node string, callbacks ...EventCallbacks) {
	if len(callbacks) == 0 {
		return
	}

	children, _, eventCh, err := z.Conn.ChildrenW(node)
	if err != nil {
		log.Printf("ChildrenW(%s): %s", node, err)
		return
	}

	for _, child := range children {
		childNode := path.Join(node, child)
		z.ManageNode(childNode, callbacks[0])
		if len(callbacks) > 1 {
			go z.ManageTree(childNode, callbacks[1:]...)
		}
	}

	for {
		ev := <-eventCh

		if ev.State == zookeeper.STATE_CLOSED {
			// shutdown was called on ZkConn?
			return
		}

		if ev.State == zookeeper.STATE_EXPIRED_SESSION ||
			ev.State == zookeeper.STATE_CONNECTING {
			log.Printf("ChildrenW(%s): exiting on loss of connection", node)
			return
		}

		switch ev.Type {
		case zookeeper.EVENT_DELETED:
			log.Printf("ChildrenW(%s): exiting on EVENT_DELETED", node)
			return
		case zookeeper.EVENT_CHILD:
			prev := children
			children, _, eventCh, err = z.Conn.ChildrenW(node)
			if err != nil {
				log.Printf("ChildrenW(%s): %s", node, err)
				return
			}
			for _, child := range ArrayDiff(children, prev) {
				childNode := path.Join(node, child)
				z.ManageNode(childNode, callbacks[0])
				if len(callbacks) > 1 {
					go z.ManageTree(childNode, callbacks[1:]...)
				}
			}
		}
	}
}
