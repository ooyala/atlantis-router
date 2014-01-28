package router

import (
	"atlantis/router/logger"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StatusServer struct {
	router *Router
}

func NewStatusServer(r *Router) *StatusServer {
	return &StatusServer{
		router: r,
	}
}

func (s *StatusServer) StatusZ(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statusz.html")
}

func (s *StatusServer) StatusZJSON(w http.ResponseWriter, r *http.Request) {
	json, err := s.router.config.StatusZJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, json)
}

func (s *StatusServer) PrintRouting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("content-type", "text/plain")

	// modify request to get rid of /port
	r.URL.Path = strings.Replace(r.URL.Path, "/"+vars["port"], "", 1)
	r.RequestURI = strings.Replace(r.RequestURI, "/"+vars["port"], "", 1)

	port, err := strconv.ParseUint(vars["port"], 10, 16)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, s.router.config.PrintRouting(uint16(port), r))
}

func (s *StatusServer) Run(port uint16, tout time.Duration) {
	gmux := mux.NewRouter()
	gmux.HandleFunc("/statusz", s.StatusZ).Methods("GET")
	gmux.HandleFunc("/statusz.json", s.StatusZJSON).Methods("GET")
	gmux.PathPrefix("/{port:[0-9]+}").HandlerFunc(s.PrintRouting)

	server := http.Server{
		Handler:      gmux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  tout,
		WriteTimeout: tout,
	}

	for {
		logger.Errorf("[status server] %s", server.ListenAndServe())
		time.Sleep(1 * time.Second)
	}
}
