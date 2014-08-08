
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)

//TODO: before calling api/zk methods, authenticate
func ListTries(w http.ResponseWriter, r *http.Request) {
	
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
	
		m := GetMapFromReqJson(r)	
		//user := m["User"]
		//secret := m["Secret"]
	} 

	tries, err := zk.ListTries()
		
}

func GetTrie(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
				
	} 

	trie, err := zk.GetTrie(vars["TrieName"])
}


func SetTrie(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")

	var trie cfg.Trie
	//Accept incoming as Json
	if contentType == "application/json" {

		body, err := GetRequestBody(r)
		if err != nil {
			//error
		}
		err = json.Unmarshal(body, &trie)		
		if err != nil {
			//return some error or something
		}

	} 

	//handle return	
	err := zk.SetTrie(trie)

}

func DeleteTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if contentType == "application/json" {
		
		m, err := GetMapFromReqJson(r)
		if err != nil {
			//error
		}	
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
	
	} 
	

	err := zk.DeleteTrie(vars["TrieName"])	
	
}
