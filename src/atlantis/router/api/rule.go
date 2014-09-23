
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)


func ListRules(w http.ResponseWriter, r *http.Request) {

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
		
	rules, err := zk.ListRules()
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
		
	var rMap map[string][]string
	rMap["Rules"] = rules

	rJson, err := json.Marshal(rMap)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, rJson)	
}

func GetRule(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}

	rule, err := zk.GetRule(vars["RuleName"])
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	if rule.Name == "" {
		WriteResponse(w, NotFoundStatusCode, GetStatusJson(ResourceDoesNotExistStatus + ": " + vars["RuleName"]))
		return
	}

	rJson, err := json.Marshal(rule)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, rJson)
}


func SetRule(w http.ResponseWriter, r *http.Request) {
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

	var rule cfg.Rule
	err = json.Unmarshal(body, &rule)		
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}

	
	err := zk.SetRule(rule)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))

}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
 

	err = zk.DeleteRule(vars["RuleName"])	
	
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}
