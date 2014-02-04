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

// The oddity that is the name of this file can be explained. Gocov will compute
// coverage on functions defined in this file if it were named test_helpers.go!
// In the interest of maintaining readable coverage statistics, and helping the
// go tools realize this file is related to testing, we let master Yoda name it.

package config

func rootTrie() Trie {
	return Trie{
		Name: "root",
		Rules: []string{
			"pastaRule",
			"meatRule",
			"cookieRule",
		},
		Internal: false,
	}
}

func pastaRule() Rule {
	return Rule{
		Name:     "pastaRule",
		Type:     "host",
		Value:    "pasta.atlantis/router.net",
		Pool:     "",
		Next:     "pastaTrie",
		Internal: false,
	}
}

func pastaTrie() Trie {
	return Trie{
		Name: "pastaTrie",
		Rules: []string{
			"fettuccinePastaRule",
			"ravioliPastaRule",
			"spaghettiPastaRule",
			"macaroniPastaRule",
		},
		Internal: false,
	}
}

func fettuccinePastaRule() Rule {
	return Rule{
		Name:     "fettuccinePastaRule",
		Type:     "path-prefix",
		Value:    "/fetuccine",
		Pool:     "pastaPool",
		Next:     "",
		Internal: false,
	}
}

func ravioliPastaRule() Rule {
	return Rule{
		Name:     "ravioliPastaRule",
		Type:     "path-prefix",
		Value:    "/ravioli",
		Pool:     "pastaPool",
		Next:     "",
		Internal: false,
	}
}

func spaghettiPastaRule() Rule {
	return Rule{
		Name:     "spaghettiPastaRule",
		Type:     "path-prefix",
		Value:    "/spaghetti",
		Pool:     "pastaPool",
		Next:     "",
		Internal: false,
	}
}

func macaroniPastaRule() Rule {
	return Rule{
		Name:     "macaroniPastaRule",
		Type:     "path-prefix",
		Value:    "/macaroni",
		Pool:     "pastaPool",
		Next:     "",
		Internal: false,
	}
}

func pastaPool() Pool {
	return Pool{
		Name:     "pastaPool",
		Internal: false,
		Hosts: map[string]Host{
			"host0": Host{
				Address: "localhost:8081",
			},
		},
		Config: PoolConfig{
			HealthzEvery:   "1s",
			HealthzTimeout: "1s",
			RequestTimeout: "1s",
			Status:         "OK",
		},
	}
}

func meatRule() Rule {
	return Rule{
		Name:     "meatRule",
		Type:     "host",
		Value:    "meat.atlantis/router.net",
		Pool:     "",
		Next:     "meatTrie",
		Internal: false,
	}
}

func meatTrie() Trie {
	return Trie{
		Name: "meatTrie",
		Rules: []string{
			"rabbitMeatRule",
			"bisonMeatRule",
			"alligatorMeatRule",
			"kangarooMeatRule",
		},
		Internal: false,
	}
}

func rabbitMeatRule() Rule {
	return Rule{
		Name:     "rabbitMeatRule",
		Type:     "path-prefix",
		Value:    "/rabbit",
		Pool:     "butcheryPool",
		Next:     "",
		Internal: false,
	}
}

func bisonMeatRule() Rule {
	return Rule{
		Name:     "bisonMeatRule",
		Type:     "path-prefix",
		Value:    "/bison",
		Pool:     "butcheryPool",
		Next:     "",
		Internal: false,
	}
}

func alligatorMeatRule() Rule {
	return Rule{
		Name:     "alligatorMeatRule",
		Type:     "path-prefix",
		Value:    "/alligator",
		Pool:     "butcheryPool",
		Next:     "",
		Internal: false,
	}
}

func kangarooMeatRule() Rule {
	return Rule{
		Name:     "kangarooMeatRule",
		Type:     "path-prefix",
		Value:    "/kangaroo",
		Pool:     "butcheryPool",
		Next:     "",
		Internal: false,
	}
}

func butcheryPool() Pool {
	return Pool{
		Name:     "butcheryPool",
		Internal: false,
		Hosts: map[string]Host{
			"host0": Host{
				Address: "localhost:8082",
			},
		},
		Config: PoolConfig{
			HealthzEvery:   "1s",
			HealthzTimeout: "1s",
			RequestTimeout: "1s",
			Status:         "OK",
		},
	}
}

func cookieRule() Rule {
	return Rule{
		Name:     "cookieRule",
		Type:     "host",
		Value:    "cookies.atlantis/router.net",
		Pool:     "",
		Next:     "cookieTrie",
		Internal: false,
	}
}

func cookieTrie() Trie {
	return Trie{
		Name: "cookieTrie",
		Rules: []string{
			"butterCookieRule",
			"sugarCookieRule",
			"gingerCookieRule",
			"oreoCookieRule",
		},
		Internal: false,
	}
}

func butterCookieRule() Rule {
	return Rule{
		Name:     "butterCookieRule",
		Type:     "path-prefix",
		Value:    "/butter",
		Pool:     "bakeryPool",
		Next:     "",
		Internal: false,
	}
}

func sugarCookieRule() Rule {
	return Rule{
		Name:     "sugarCookieRule",
		Type:     "path-prefix",
		Value:    "/sugar",
		Pool:     "bakeryPool",
		Next:     "",
		Internal: false,
	}

}

func gingerCookieRule() Rule {
	return Rule{
		Name:     "gingerCookieRule",
		Type:     "path-prefix",
		Value:    "/ginger",
		Pool:     "bakeryPool",
		Next:     "",
		Internal: false,
	}
}

func oreoCookieRule() Rule {
	return Rule{
		Name:     "oreoCookieRule",
		Type:     "path-prefix",
		Value:    "/oreo",
		Pool:     "nabiscoPool",
		Next:     "",
		Internal: false,
	}
}

func bakeryPool() Pool {
	return Pool{
		Name:     "bakeryPool",
		Internal: false,
		Hosts: map[string]Host{
			"host0": Host{
				Address: "localhost:8083",
			},
		},
		Config: PoolConfig{
			HealthzEvery:   "1s",
			HealthzTimeout: "1s",
			RequestTimeout: "1s",
			Status:         "OK",
		},
	}
}

func nabiscoPool() Pool {
	return Pool{
		Name:     "nabiscoPool",
		Internal: false,
		Hosts: map[string]Host{
			"host0": Host{
				Address: "localhost:8084",
			},
		},
		Config: PoolConfig{
			HealthzEvery:   "1s",
			HealthzTimeout: "1s",
			RequestTimeout: "1s",
			Status:         "OK",
		},
	}
}

func meatPort() Port {
	return Port{
		Port:     uint16(8080),
		Trie:     "meatTrie",
		Internal: false,
	}
}

func rootPort() Port {
	return Port{
		Port:     uint16(8081),
		Trie:     "rootTrie",
		Internal: false,
	}
}
