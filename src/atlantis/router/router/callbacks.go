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

package router

import (
	"atlantis/router/config"
	"atlantis/router/logger"
	"atlantis/router/zk"
	"encoding/json"
	"path"
	"strconv"
)

type PoolCallbacks struct {
	config *config.Config
}

func (p *PoolCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("PoolCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		logger.Errorf("%s unmarshalling %s as pool", err.Error(), jsonBlob)
		return
	}
	p.config.AddPool(zkPool.Pool(map[string]config.Host{}))
}

func (p *PoolCallbacks) Deleted(zkPath string) {
	logger.Debugf("PoolCallbacks.Deleted(%s)", zkPath)
	p.config.DelPool(path.Base(zkPath))
}

func (p *PoolCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("PoolCallbacks.Changed(%s, %s)", path, jsonBlob)
	var zkPool zk.ZkPool
	if err := json.Unmarshal([]byte(jsonBlob), &zkPool); err != nil {
		logger.Errorf("%s unmarshalling %s as pool", err.Error(), jsonBlob)
		return
	}
	p.config.UpdatePool(zkPool.Pool(nil))
}

type HostCallbacks struct {
	config *config.Config
}

func (h *HostCallbacks) splitPath(zkPath string) (string, string) {
	return path.Base(zkPath), path.Base(path.Dir(zkPath))
}

func (h *HostCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("HostCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	hostName, poolName := h.splitPath(zkPath)

	var host config.Host
	if err := json.Unmarshal([]byte(jsonBlob), &host); err != nil {
		logger.Errorf("%s unmarshalling %s as host", err.Error(), jsonBlob)
		return
	}

	if pool := h.config.Pools[poolName]; pool != nil {
		pool.AddServer(hostName, h.config.ConstructServer(host))
	}
}

func (h *HostCallbacks) Deleted(zkPath string) {
	logger.Debugf("HostCallbacks.Deleted(%s)", zkPath)
	hostName, poolName := h.splitPath(zkPath)
	if pool := h.config.Pools[poolName]; pool != nil {
		pool.DelServer(hostName)
	}
}

func (h *HostCallbacks) Changed(path, jsonBlob string) {
	logger.Errorf("HostCallbacks.Changed(%s, %s)", path, jsonBlob)
}

type RuleCallbacks struct {
	config *config.Config
}

func (p *RuleCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("RuleCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		logger.Errorf("%s unmarshalling %s as rule", err.Error(), jsonBlob)
		return
	}
	p.config.AddRule(rule)
}

func (p *RuleCallbacks) Deleted(zkPath string) {
	logger.Debugf("RuleCallbacks.Deleted(%s)", zkPath)
	p.config.DelRule(path.Base(zkPath))
}

func (p *RuleCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("RuleCallbacks.Changed(%s, %s)", path, jsonBlob)
	var rule config.Rule
	if err := json.Unmarshal([]byte(jsonBlob), &rule); err != nil {
		logger.Errorf("%s unmarshalling %s as rule", err.Error(), jsonBlob)
		return
	}
	p.config.UpdateRule(rule)
}

type TrieCallbacks struct {
	config *config.Config
}

func (p *TrieCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("TrieCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var trie config.Trie
	err := json.Unmarshal([]byte(jsonBlob), &trie)
	if err != nil {
		logger.Errorf("%s unmarshalling %s as trie", err.Error(), jsonBlob)
		return
	}
	p.config.AddTrie(trie)
}

func (p *TrieCallbacks) Deleted(zkPath string) {
	logger.Debugf("TrieCallbacks.Deleted(%s)", zkPath)
	p.config.DelTrie(path.Base(zkPath))
}

func (p *TrieCallbacks) Changed(path, jsonBlob string) {
	logger.Debugf("TrieCallbacks.Changed(%s, %s)", path, jsonBlob)
	var trie config.Trie
	if err := json.Unmarshal([]byte(jsonBlob), &trie); err != nil {
		logger.Errorf("%s unmarshalling %s as trie", err.Error(), jsonBlob)
		return
	}
	p.config.UpdateTrie(trie)
}

type PortCallbacks struct {
	config *config.Config
	router *Router
}

func (p *PortCallbacks) Created(zkPath, jsonBlob string) {
	logger.Debugf("PortCallbacks.Created(%s, %s)", zkPath, jsonBlob)
	var port config.Port
	if err := json.Unmarshal([]byte(jsonBlob), &port); err != nil {
		logger.Errorf("%s unmarshalling %s as port", err.Error(), jsonBlob)
		return
	}
	p.config.AddPort(port)
	p.router.AddPort(port.Port)
}

func (p *PortCallbacks) Deleted(zkPath string) {
	logger.Debugf("PortCallbacks.Deleted(%s)", zkPath)
	port, err := strconv.ParseUint(path.Base(zkPath), 10, 16)
	if err != nil {
		logger.Errorf("%s interpreting base of %s as uint16", err.Error(), zkPath)
		return
	}
	p.config.DelPort(uint16(port))
	p.router.DelPort(uint16(port))
}

func (p *PortCallbacks) Changed(zkPath, jsonBlob string) {
	logger.Debugf("PortCallbacks.Changed(%s)", zkPath, jsonBlob)
	var port config.Port
	if err := json.Unmarshal([]byte(jsonBlob), &port); err != nil {
		logger.Errorf("%s unmarshalling %s as port", err.Error(), jsonBlob)
		return
	}
	p.config.UpdatePort(port)
}
