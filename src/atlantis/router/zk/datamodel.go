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
	"atlantis/router/config"
)

type ZkPool struct {
	Name     string
	Internal bool
	Config   config.PoolConfig
}

func ToZkPool(p config.Pool) (ZkPool, map[string]config.Host) {
	zkPool := ZkPool{
		Name:     p.Name,
		Internal: p.Internal,
		Config:   p.Config,
	}

	return zkPool, p.Hosts
}

func (z ZkPool) Pool(hosts map[string]config.Host) config.Pool {
	return config.Pool{
		Name:     z.Name,
		Internal: z.Internal,
		Hosts:    hosts,
		Config:   z.Config,
	}
}
