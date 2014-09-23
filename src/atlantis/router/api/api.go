


package api


import (
	"fmt"
	"net"
	"net/http"
	"github.com/gorilla/mux"
)

const (
// Status Strings and a colon which will seperate them from specific reason

	NotAuthenticatedStatus = "Authentication not succesful: "
	IncorrectContentTypeStatus = "Request failed, the content-type must be application/json"
	CouldNotReadRequestDataStatus = "Failed to read request data: "
	CouldNotCompleteOperationStatus = "The requested operation could not be completed: "
	ResourceDoesNotExistStatus = "The requested resources does not exist"
	RequestSuccesfulStatus = "The request was succesfull"
)

const (
// Status codes supported for the API

	OkStatusCode = 200
	BadRequestStatusCode = 400
	NotAuthorizedStatusCode = 401
	NotFoundStatusCode = 404
	ServerErrorCode = 500
) 




func NotFound(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, notFoundHTML)
}


func Init(listenAddr string) error{

	gmux = mux.NewRouter()
	
	gmux.NotFoundHandler = http.HandlerFunc(NotFound)

	/*	
	//router management
	gmux.HandleFunc("/routers", ListRouters).Methods("GET")
	gmux.HandleFunc("/routers/zone/{Zone}/host/{Host}", GetRouter).Methods("GET")
	gmux.HandleFunc("/routers/zone/{Zone}/host/{Host}", RegisterRouter).Methods("PUT")
	gmux.HandleFunc("/routers/zone/{Zone}/host/{Host}", UnregisterRouter).Methods("DELETE")
	*/

	//router config
	
	//Pools
	gmux.HandleFunc("/pools", ListPools).Methods("GET")
	gmux.HandleFunc("/pools/{PoolName}", GetPool).Methods("GET")
	gmux.HandleFunc("/pools/{PoolName}", SetPool).Methods("PUT")
	gmux.HandleFunc("/pools/{PoolName}", DeletePool).Methods("DELETE")	
	
	//hosts for pool
	gmux.HandleFunc("/pools/{PoolName}/hosts", GetHosts).Methods("GET")
	gmux.HandleFunc("/pools/{PoolName}/hosts", AddHosts).Methods("PUT")
	gmux.HandleFunc("/pools/{PoolName}/hosts", DeleteHosts).Methods("DELETE")



	//Rules
	gmux.HandleFunc("/rules", ListRules).Methods("GET")
	gmux.HandleFunc("/rules/{RuleName}", GetRule).Methods("GET")
	gmux.HandleFunc("/rules/{RuleName}", SetRule).Methods("PUT")
	gmux.HandleFunc("/rules/{RuleName}", DeleteRule).Methods("DELETE")
	
	//Tries
	gmux.HandleFunc("/tries", ListTries).Methods("GET")
	gmux.HandleFunc("/tries/{TrieName}", GetTrie).Methods("GET")
	gmux.HandleFunc("/tries/{TrieName}", SetTrie).Methods("PUT")
	gmux.HandleFunc("/tries/{TrieName}", DeleteTrie).Methods("DELETE")
	
	//Ports
	//gmux.HandleFunc("/ports/apps/{App}/envs/{Env}", GetAppEnvPort).Methods("GET")
	//gmux.HandleFunc("/ports/apps", ListAppEnvsWithPort).Methods("GET")

	gmux.HandleFunc("/ports", ListPorts).Methods("GET")	
	gmux.HandleFunc("/ports/{Port}", GetPort).Methods("GET")
	gmux.HandleFunc("/ports/{Port}", SetPort).Methods("PUT")
	gmux.HandleFunc("/ports/{Port}", DeletePort).Methods("DELETE")

	handler := apachelog.NewHandler(HandlerFunc(gmux), os.Stderr)
	server = &http.Server{Addr: listenAddr, Handler: handler}
	lAddr = listenAddr
	return nil

}

func Listen() {

	if server == nil {
		panic("Not Initialized.")
	}

}
