
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)

//TODO: before calling api/zk methods, authenticate
func ListPorts(w http.ResponseWriter, r *http.Request) {
	
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
	
		m := GetMapFromReqJson(r)	
		//user := m["User"]
		//secret := m["Secret"]
	} 

	ports, err := zk.ListPorts()
		
}

func GetPort(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
				
	}
 
	port, err := zk.GetPort(vars["PortName"])
}


func SetPort(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")

	var port cfg.Port
	//Accept incoming as Json
	if contentType == "application/json" {

		body, err := GetRequestBody(r)
		if err != nil {
			//error
		}
		err = json.Unmarshal(body, &port)
		if err != nil {
			//return some error or something
		}

	} 

	//handle return	
	err := zk.SetPort(port)

}

func DeletePort(w http.ResponseWriter, r *http.Request) {
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
	
	err = zk.DeletePort(vars["PortName"])	
	
}
