
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)

func ListPorts(w http.ResponseWriter, r *http.Request) {

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}	
	
	ports, err := zk.ListPorts()
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	var pMap map[string][]uint16
	pMap["Ports"] = ports		
	
	pJson, err := json.Marshal(pMap)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, pJson)
	
}

func GetPort(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
	 
	port, err := zk.GetPort(vars["PortName"])
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	if port.Port == 0 {
		WriteResponse(w, NotFoundStatusCode, GetStatusJson(ResourceDoesNotExistStatus + ": " + vars["PortName"]))
		return
	}

	pJson, err := json.Marshal(port)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, pJson)
}


func SetPort(w http.ResponseWriter, r *http.Request) {

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
	
	//Accept incoming as Json
	if r.Header.Get("Content-Type") != "application/json" {
		WriteResponse(w, BadRequestStatusCode, GetStatusJson(IncorrectContentTypeStatus))
		return
	}
	
	body, err := GetRequestBody(r)
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}
	var port cfg.Port
	err = json.Unmarshal(body, &port)
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}

	
	err = zk.SetPort(port)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}

func DeletePort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
	
	err = zk.DeletePort(vars["PortName"])	
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}
