
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)

//TODO: before calling api/zk methods, authenticate
func ListRouters(w http.ResponseWriter, r *http.Request) {
	
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
	
		m := GetMapFromReqJson(r)	
		//user := m["User"]
		//secret := m["Secret"]
	} 

	routers, err := zk.ListRouters()
		
}

func RegisterRouter(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
		//vars["Zone"]
		//vars["Host"]
		//ip = m["IP"}
				
	} 

	id, err := zk.RegisterRouter()
}

func UnregisterRouter(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {

		m := GetMapFromReqJson(r) 
		//name := m["Name"]
		//user := m["User"]
		//secret := m["Secret"]
		//vars["Zone"]
		//vars["Host"]
		//ip = m["IP"]
				
	} 

	id, err := zk.UnregisterRouter()
}

func GetRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {
		m := GetMapFromReqJson(r)
		
		//user := m["User"]
		//secret := m["Secret"]
		//vars["Zone"]
		//vars["Host"]
		//ip = m["IP"]

	}

	router, err := zk.GetRouter()	

}
