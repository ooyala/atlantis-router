
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)

//TODO: before calling api/zk methods, authenticate
func ListRules(w http.ResponseWriter, r *http.Request) {
	
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
	
		m := GetMapFromReqJson(r)	
		//user := m["User"]
		//secret := m["Secret"]
	} 

	rules, err := zk.ListRules()
		
}

func GetRule(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
				
	} 

	rule, err := zk.GetRule(vars["RuleName"])
}


func SetRule(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")

	var rule cfg.Rule
	//Accept incoming as Json
	if contentType == "application/json" {

		body, err := GetRequestBody(r)
		if err != nil {
			//error
		}
		err = json.Unmarshal(body, &rule)		
		if err != nil {
			//return some error or something
		}

	} 

	//handle return	
	err := zk.SetRule(rule)

}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
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
	

	err := zk.DeleteRule(vars["RuleName"])	
	
}
