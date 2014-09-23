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
	
	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
	 
	hostsMap, err := zk.GetHosts(vars["PoolName"])
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	hMapJson, err := json.Marshal(hostsMap)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, hMapJson)
}


func AddHosts(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
		
	if r.Header.Get("Content-Type") != "application/json" {
		WriteResponse(w, BadRequestStatusCode, GetStatusJson(IncorrectContentTypeStatus))
		return
	}

	var hostMap map[string]cfg.Host
	body, err := GetRequestBody(r)
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}
	err = json.Unmarshal(body, &hostsMap)		
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))	
		return
	}

	err = zk.AddHosts(vars["PoolName"], hostsMap)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))

}

func DeleteHosts(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
                return
	}
	
	if r.Header.Get("Content-Type") != "application/json" {
		WriteResponse(w, BadRequestStatusCode, GetStatusJson(IncorrectContentTypeStatus))
                return
	}

	m := GetMapFromReqJson(r) 

	var hostList []string	
	hList := m["Hosts"]
	fList := hList.([]interface{})

	//parse the standard host req format to adjust for rw.go format
	for key, value := range fList {
		temp := value.(map[string]interface{})
		hostList[key] = temp["Address"]	   
	}

	err = zk.DeleteHosts(vars["PoolName"], hostList)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}

