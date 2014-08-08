
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"atlantis/router/api/auth"
	"atlantis/router/api/types"
	cfg "atlantis/router/config"
)

//TODO: before calling api/zk methods, authenticate
func ListPools(w http.ResponseWriter, r *http.Request) {
	
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		//error	
	}	
	
	m := GetMapFromReqJson(r)	
	user := m["User"]
	secret := m["Secret"]

	err := auth.IsAllowed(user, secret)
	if err != nil {
		//not auth error/return
	}

	pools, err := zk.ListPools()
	if err != nil {
		//error
 	}

 	poolsJson, err := json.Marshal(pools)
	if err != nil {
		//error 
	}
	
	fmt.Fprintf(w, "%s", poolsJson)	
		
}

func GetPool(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")
	
	if contentType != "application/json" {
		//error
	}

	m := GetMapFromReqJson(r) 
	user := m["User"]
	secret := m["Secret"]
				
	err := auth.IsAllowed(user, secret)
	if err != nil {
		//not auth/exit
	}	 

	pool, err := zk.GetPool(vars["PoolName"])
	if err != nil {
		//error
	}

   	poolJson, err := json.Marshal(pool)
	if err != nil {
		//error
	}	
	
	fmt.Fprintf(w, "%s", poolJson)
}


func SetPool(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		//error exit
	}

	b, err := GetRequestBody(r)
	if err != nil {
		//error
	}

	var apiPool types.ApiPool 
	err = json.Unmarshal(body, &apiPool)		
	if err != nil {
		//return some error or something
	}

	err := auth.IsAllowed(apiPool.User, apiPool.Secret)


	err := zk.SetPool(apiPool.Pool)

	if err != nil {
		//error
	}

	fmt.Fprintf(w, "%s", GetStatusJson("Pool set succesfully")) 

}

func DeletePool(w http.ResponseWriter, r *http.Request) {
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
	

	err := zk.DeletePool(vars["PoolName"])	

	if err != nil {
		//error
	}

	fmt.Fprintf(w, "%s", GetStatusJson(vars["PoolName"] + " deleted succesfully"))	
}
