
package api


import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	cfg "atlantis/router/config"
)


func ListTries(w http.ResponseWriter, r *http.Request) {

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}	

	tries, err := zk.ListTries()
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	var tMap map[string][]string
	tMap["Tries"] = tries
	
	tJson, err := json.Marshal(tMap)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, tJson)
		
}

func GetTrie(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}	

	trie, err := zk.GetTrie(vars["TrieName"])
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	if trie.Name == "" {
		WriteResponse(w, NotFoundStatusCode, GetStatusJson(ResourceDoesNotExistStatus + ": " + vars["TrieName"]))
		return
	} 

	tJson, err := json.Marshal(trie)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, tJson)
}


func SetTrie(w http.ResponseWriter, r *http.Request) {

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
	
	var trie cfg.Trie
	err = json.Unmarshal(body, &trie)		
	if err != nil {
		WriteResponse(w, BadRequestStatusCode, GetErrorStatusJson(CouldNotReadRequestDataStatus, err))
		return
	}

	
	err = zk.SetTrie(trie)
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}

	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}

func DeleteTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := GetUserSecretAndAuth(r)
	if err != nil {
		WriteResponse(w, NotAuthorizedStatusCode, GetErrorStatusJson(NotAuthenticatedStatus, err))
		return
	}
	
	err := zk.DeleteTrie(vars["TrieName"])	
	if err != nil {
		WriteResponse(w, ServerErrorCode, GetErrorStatusJson(CouldNotCompleteOperationStatus, err))
		return
	}
	
	WriteResponse(w, OkStatusCode, GetStatusJson(RequestSuccesfulStatus))
}
