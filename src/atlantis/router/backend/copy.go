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
	"container/list"
	"io"
	"sync"
	"time"
)

// Profiling shows considerable amount of time spent in io.Copy() spent trying to allocate memory. The reason for
// poor performance being that Linux cannot use sendfile(3) to copy between sockets. Consequently, io.Copy() uses
// GenericReadFrom on io.Reader and GenericWriteTo on io.Writer, which internally allocate a 32kB buffer on every
// copy, creating memory pressure for GC and considerable time spent in gc.malloc(). As a work around, we imlpement
// our own copier, which re-uses buffers. The Copy() function is a near replica of io.Copy() in the generic case.

type queued struct {
	when time.Time
	buf  []byte
}

type Copier struct {
	sync.Mutex
	queue *list.List
	tickF time.Duration
	killC chan bool
}

func NewCopier(args ...time.Duration) *Copier {
	copier := &Copier{
		queue: new(list.List),
		tickF: 1 * time.Minute,
		killC: make(chan bool),
	}
	if args != nil {
		copier.tickF = args[0]
	}

	go copier.metronome()
	return copier
}

func (c *Copier) metronome() {
	tick := time.Tick(c.tickF)

	for {
		select {
		case <-tick:
			c.Lock()
			c.gcList()
			c.Unlock()
		case <-c.killC:
			return
		}
	}
}

func (c *Copier) gcList() {
	// Low threshold
	if c.queue.Len() < 8 {
		return
	}

	for e := c.queue.Front(); e != nil; e = e.Next() {
		if time.Since(e.Value.(queued).when) > c.tickF {
			c.queue.Remove(e)
			e.Value = nil
		}
	}
}

func (c *Copier) getBuffer() []byte {
	c.Lock()
	defer c.Unlock()

	if c.queue.Len() == 0 {
		return make([]byte, 32*1024)
	}

	e := c.queue.Front()
	c.queue.Remove(e)
	return e.Value.(queued).buf
}

func (c *Copier) putBuffer(buf []byte) {
	c.Lock()
	defer c.Unlock()

	c.queue.PushFront(queued{when: time.Now(), buf: buf})
}

func (c *Copier) Copy(dst io.Writer, src io.Reader) (wr int64, err error) {
	buf := c.getBuffer()
	defer c.putBuffer(buf)

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				wr += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return wr, err
}

func (c *Copier) Shutdown() {
	c.Lock()
	defer c.Unlock()

	for e := c.queue.Front(); e != nil; e = e.Next() {
		c.queue.Remove(e)
		e.Value = nil
	}

	c.killC <- true
}
