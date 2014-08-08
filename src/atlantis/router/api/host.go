
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
	"atlantis/router/api/auth"
)

func GetHosts(w http.ResponseWriter, r *http.Request) {

	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		user := m["User"]
		secret := m["Secret"]
		isAuth, err := auth.SimpleAuth(user, secret)
		if err != nil {
			//error authenticating
		}		
	 
	if isAuth { 
		hostsMap, err := zk.GetHosts(vars["PoolName"])
	}else {
		//not auth
	}
}


func AddHosts(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")

	var hostsMap map[string]cfg.Host

	//Accept incoming as Json
	if contentType == "application/json" {

		body, err := GetRequestBody(r)
		if err != nil {
			//error
		}
		err = json.Unmarshal(body, &hostsMap)		
		if err != nil {
			//return some error or something
		}

	} 

	err := zk.AddHosts(vars["PoolName"], hostsMap)

}

func DeleteHosts(w http.ResponseWriter, r *http.Request) {

	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	var hostList []string	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
	
		hList := m["Hosts"]
		fList := hList.([]interface{})

		for key, value := range fList {
			temp := value.(map[string]interface{})
			hostList[key] = temp["Address"]	   
		}
	} 

	err := zk.DeleteHosts(vars["PoolName"], hostList)
}

