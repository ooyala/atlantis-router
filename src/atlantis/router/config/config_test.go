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
	"atlantis/router/routing"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())
	if config.MatcherFactory == nil {
		t.Errorf("should set matcher factory")
	}
	if config.Pools == nil || config.Rules == nil || config.Tries == nil {
		t.Errorf("should make maps")
	}
}

func TestAddPool(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	config.AddRule(butterCookieRule())
	config.AddRule(sugarCookieRule())
	config.AddRule(gingerCookieRule())
	config.AddRule(oreoCookieRule())

	config.AddPool(bakeryPool())
	defer config.DelPool("bakeryPool")

	if config.Rules["butterCookieRule"].PoolPtr != config.Pools["bakeryPool"] ||
		config.Rules["sugarCookieRule"].PoolPtr != config.Pools["bakeryPool"] ||
		config.Rules["gingerCookieRule"].PoolPtr != config.Pools["bakeryPool"] {
		t.Errorf("should update references to bakery pool")
	}

	if config.Rules["oreoCookieRule"] == config.Rules["butterCookieRule"] {
		t.Errorf("shoud not update references to nabisco pool")
	}

	current := config.Pools["bakeryPool"]
	config.AddPool(bakeryPool())
	if current != config.Pools["bakeryPool"] {
		t.Errorf("should silently ignore existing pools")
	}
}

func TestUpdatePool(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent pools")
		}
	}()
	config.UpdatePool(bakeryPool())

	config.AddPool(bakeryPool())
	config.AddRule(butterCookieRule())
	config.AddRule(sugarCookieRule())

	pool := bakeryPool()
	conf := PoolConfig{
		HealthzEvery:   "60s",
		HealthzTimeout: "60s",
		RequestTimeout: "60s",
		Status:         "CRITICAL",
	}
	pool.Config = conf
	config.UpdatePool(pool)
	if config.Pools["bakeryPool"].Config.Status != "CRITICAL" {
		t.Errorf("should update pool")
	}

	if config.Rules["butterCookieRule"].PoolPtr != config.Pools["bakeryPool"] ||
		config.Rules["sugarCookieRule"].PoolPtr != config.Pools["bakeryPool"] {
		t.Errorf("should update references to bakery pool")
	}
}

func TestDelPool(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent pools")
		}
	}()
	config.DelPool("bakeryPool")

	config.AddPool(bakeryPool())
	config.AddPool(nabiscoPool())

	config.AddRule(butterCookieRule())
	config.AddRule(sugarCookieRule())
	config.AddRule(gingerCookieRule())
	config.AddRule(oreoCookieRule())

	config.AddPool(pastaPool())
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently delete pools without references")
		}
	}()
	config.DelPool("pastaPool")

	config.DelPool("bakeryPool")
	if _, ok := config.Pools["bakeryPool"]; ok {
		t.Errorf("should delete bakery pool")
	}
	if config.Rules["butterCookieRule"].PoolPtr != nil ||
		config.Rules["sugarCookieRule"].PoolPtr != nil ||
		config.Rules["gingerCookieRule"].PoolPtr != nil {
		t.Errorf("should nil references to bakery pool")
	}

	if config.Rules["oreoCookieRule"].PoolPtr == nil {
		t.Errorf("should not nil reference to nabisco pool")
	}
}

func TestAddRule(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently construct rules missing tries/ pools")
		}
	}()
	config.AddRule(pastaRule())
	config.AddTrie(pastaTrie())

	config.AddRule(fettuccinePastaRule())
	if _, ok := config.Rules["fettuccinePastaRule"]; !ok {
		t.Errorf("should add fettuccine pasta rule %#v", ok)
	}

	config.AddRule(ravioliPastaRule())
	if config.Tries["pastaTrie"].List[0] != config.Rules["fettuccinePastaRule"] ||
		config.Tries["pastaTrie"].List[1] != config.Rules["ravioliPastaRule"] {
		t.Errorf("should update references to fettuccine and spaghetti rule")
	}

	fettuccine := config.Rules["fettuccinePastaRule"]
	config.AddRule(fettuccinePastaRule())
	if config.Rules["fettuccinePastaRule"] != fettuccine {
		t.Errorf("should silently ignore existing rules")
	}
}

func TestUpdateRule(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent rules")
		}
	}()
	config.UpdateRule(pastaRule())

	config.AddPool(pastaPool())
	config.AddTrie(pastaTrie())
	config.AddRule(fettuccinePastaRule())

	config.AddPool(nabiscoPool())
	rule := fettuccinePastaRule()
	rule.Pool = "nabiscoPool"
	config.UpdateRule(rule)
	if config.Rules["fettuccinePastaRule"].PoolPtr != config.Pools["nabiscoPool"] {
		t.Errorf("should update rule")
	}

	if config.Tries["pastaTrie"].List[0] != config.Rules["fettuccinePastaRule"] {
		t.Errorf("should update references to pool")
	}
}

func TestDelRule(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent rules")
		}
	}()
	config.DelRule("pastaRule")

	config.AddPool(bakeryPool())
	config.AddRule(fettuccinePastaRule())
	config.AddRule(spaghettiPastaRule())
	config.AddTrie(pastaTrie())

	config.AddRule(kangarooMeatRule())
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently delete rules without references")
		}
	}()
	config.DelRule("kangarooMeatRule")

	fettuccine := config.Rules["fettuccinePastaRule"]
	spaghetti := config.Rules["spaghettiPastaRule"]
	config.DelRule("fettuccinePastaRule")
	if _, ok := config.Rules["fettuccinePastaRule"]; ok {
		t.Errorf("should delete fettuccine pasta rule")
	}
	if config.Tries["pastaTrie"].List[0] == fettuccine {
		t.Errorf("should update references to fettuccine pasta rule")
	}
	if config.Tries["pastaTrie"].List[2] != spaghetti {
		t.Errorf("should not modify reference to spaghetti rule")
	}
}

func TestAddTrie(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently add tries missing rules")
		}
	}()
	config.AddTrie(rootTrie())
	config.AddRule(meatRule())
	config.AddPort(meatPort())

	config.AddPool(butcheryPool())
	config.AddRule(rabbitMeatRule())
	config.AddRule(bisonMeatRule())
	config.AddRule(alligatorMeatRule())
	config.AddRule(kangarooMeatRule())
	config.AddTrie(meatTrie())

	if _, ok := config.Tries["meatTrie"]; !ok {
		t.Errorf("should add meat trie")
	}

	if config.Tries["meatTrie"].List[0] != config.Rules["rabbitMeatRule"] ||
		config.Tries["meatTrie"].List[1] != config.Rules["bisonMeatRule"] ||
		config.Tries["meatTrie"].List[2] != config.Rules["alligatorMeatRule"] ||
		config.Tries["meatTrie"].List[3] != config.Rules["kangarooMeatRule"] {
		t.Errorf("should add rules to trie")
	}

	if config.Tries["root"].List[1].NextPtr != config.Tries["meatTrie"] {
		t.Errorf("should update references to meat trie")
	}

	if config.Ports[8080] != config.Tries["meatTrie"] {
		t.Errorf("should update references to meat trie")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore existing trie")
		}
	}()
	config.AddTrie(meatTrie())
}

func TestUpdateTrie(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent rules")
		}
	}()
	config.UpdateTrie(meatTrie())
	config.AddPort(meatPort())

	config.AddPool(butcheryPool())
	config.AddRule(rabbitMeatRule())
	config.AddRule(bisonMeatRule())
	config.AddRule(alligatorMeatRule())
	config.AddRule(kangarooMeatRule())
	config.AddTrie(meatTrie())
	config.AddRule(meatRule())
	config.AddTrie(rootTrie())

	meat := meatTrie()
	meat.Rules = []string{
		"kangarooMeatRule",
		"alligatorMeatRule",
		"bisonMeatRule",
		"rabbitMeatRule",
	}

	config.UpdateTrie(meat)
	if config.Tries["meatTrie"].List[0] != config.Rules["kangarooMeatRule"] {
		t.Errorf("should update meat trie")
	}

	if config.Rules["meatRule"].NextPtr != config.Tries["meatTrie"] {
		t.Errorf("should update references to meat trie")
	}

	if config.Ports[8080] != config.Tries["meatTrie"] {
		t.Errorf("should update references to meat trie")
	}

}

func TestDelTrie(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent rules")
		}
	}()
	config.DelTrie("meatTrie")

	config.AddTrie(meatTrie())
	config.AddRule(meatRule())
	config.AddTrie(rootTrie())
	config.AddPort(meatPort())

	meatTrie := config.Tries["meatTrie"]
	config.DelTrie("meatTrie")
	if _, ok := config.Tries["meatTrie"]; ok {
		t.Errorf("should delete meat trie")
	}
	if config.Rules["meatRule"].NextPtr == meatTrie {
		t.Errorf("should nil references to meat trie")
	}
	if config.Ports[8080] == meatTrie {
		t.Errorf("should nil references to meat port")
	}
}

func TestAddPort(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently construct ports with missing tries")
		}
	}()
	config.AddPort(meatPort())

	config.AddTrie(meatTrie())
	config.AddPort(meatPort())

	if config.Ports[8080] != config.Tries["meatTrie"] {
		t.Errorf("should make port")
	}
}

func TestUpdatePort(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	config.AddTrie(rootTrie())
	config.AddTrie(meatTrie())

	root := rootPort()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently add non existent ports")
		}
	}()
	config.UpdatePort(root)

	root.Trie = "meatTrie"
	config.UpdatePort(root)
	if config.Ports[8081] != config.Tries["meatTrie"] {
		t.Errorf("should upate root port")
	}
}

func TestDelPort(t *testing.T) {
	config := NewConfig(routing.DefaultMatcherFactory())

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should silently ignore non existent ports")
		}
	}()
	config.DelPort(uint16(8080))

	config.AddTrie(rootTrie())
	config.AddPort(rootPort())
	config.AddTrie(meatTrie())
	config.AddPort(meatPort())

	config.DelPort(8081)
	if _, ok := config.Ports[8081]; ok {
		t.Errorf("should delete 8081")
	}

	if _, ok := config.Ports[8080]; !ok {
		t.Errorf("should not delete 8080")
	}
}
