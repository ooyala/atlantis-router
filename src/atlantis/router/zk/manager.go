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
	"github.com/scalingdata/gozk"
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
		logger.Errorf("[zkconn %d] GetW(%s): %s", z, node, err)
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
				logger.Printf("[zkconn %d] connection lost, stop watching %s", z, node)
				return
			}

			switch ev.Type {
			case zookeeper.EVENT_DELETED:
				callbacks.Deleted(node)
				return
			case zookeeper.EVENT_CHANGED:
				content, _, eventCh, err = z.Conn.GetW(node)
				if err != nil {
					logger.Errorf("[zkconn %d] GetW(%s): %s", z, node, err)
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
		logger.Errorf("[zkconn %d] ChildrenW(%s): %s", z, node, err)
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
			logger.Printf("[zkconn %d] connection lost, stop watching %s", z, node)
			return
		}

		switch ev.Type {
		case zookeeper.EVENT_DELETED:
			logger.Printf("[zkconn %d] node deleted, stop watching %s", z, node)
			return
		case zookeeper.EVENT_CHILD:
			prev := children
			children, _, eventCh, err = z.Conn.ChildrenW(node)
			if err != nil {
				logger.Errorf("[zkconn %d] ChildrenW(%s): %s", z, node, err)
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
