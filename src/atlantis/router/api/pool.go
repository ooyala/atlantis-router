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

	err := GetUserSecretAndAuth(r) 
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}

	pools, err := zk.ListPools()
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
 	}

 	poolsJson, err := json.Marshal(pools)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, poolsJson)	
}

func GetPool(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err :=	GetUserSecretAndAuth(r) 
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}	 

	pool, err := zk.GetPool(vars["PoolName"])
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	if pool.Name == "" {
		WriteResponse(w, NotFoundStatusCode, GetStatusJson(ResourceDoesNotExistStatus + ": " + vars["PoolName"]))
		return
	} 

   	poolJson, err := json.Marshal(pool)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	

	WriteResponse(w, OkStatusCode, poolJson)	
}


func SetPool(w http.ResponseWriter, r *http.Request) {
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

	body, err := GetRequestBody(r)
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}

	var pool cfg.Pool 
	err = json.Unmarshal(body, &pool)		
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}


	err = zk.SetPool(pool)

	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}

func DeletePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}

	err = zk.DeletePool(vars["PoolName"])	

	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}
