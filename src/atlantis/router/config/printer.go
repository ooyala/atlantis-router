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

package config

import (
	"atlantis/router/backend"
	"atlantis/router/routing"
	"fmt"
	"net/http"
)

// This leaks the abstractions of routing.Trie.Walk() and config.Route()
// and is strictly a debugging aid.
func (c *Config) PrintRouting(port uint16, r *http.Request) string {
	c.RLock()
	defer c.RUnlock()

	var next *routing.Trie
	var pool *backend.Pool

	output := fmt.Sprintf("port %d\n", port)
	indent := "  "

	next = c.Ports[port]
	for next != nil || pool != nil {
		if pool != nil {
			output += fmt.Sprintf("%spool %s\n", indent, pool.Name)
			return output
		} else {
			output += fmt.Sprintf("%strie %s\n", indent, next.Name)
		}
		indent += "  "

		for _, rule := range next.List {
			if rule.Dummy {
				output += fmt.Sprintf("%srule %s dummy\n", indent, rule.Name)
			} else if rule.Matcher.Match(r) {
				output += fmt.Sprintf("%srule %s T\n", indent, rule.Name)
				pool, next = rule.PoolPtr, rule.NextPtr
				break
			} else {
				output += fmt.Sprintf("%srule %s F\n", indent, rule.Name)
				pool, next = nil, nil
			}
		}
	}
	fmt.Sprintf("%snext = nil, pool = nil!\n", indent)
	return output
}
