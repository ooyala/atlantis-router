package router

import (
	"atlantis/router/logger"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	w.Header().Add("content-type", "application/plain")

}

func (s *StatusServer) Run(port uint16, tout time.Duration) {
	gmux := mux.NewRouter()
	gmux.HandleFunc("/statusz", s.StatusZ).Methods("GET")
	gmux.HandleFunc("/statusz.json", s.StatusZJSON).Methods("GET")
	gmux.HandleFunc("/{port:[0-9]+}/", s.PrintRouting).Methods("GET")

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
