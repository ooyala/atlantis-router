package router

import (
	"atlantis/router/logger"
	"fmt"
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

func (s *StatusServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.config.Printer(w, r)
}

func (s *StatusServer) Run(port uint16, tout time.Duration) {
	server := http.Server{
		Handler:      s,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  tout,
		WriteTimeout: tout,
	}

	for {
		logger.Errorf("[status server] %s", server.ListenAndServe())
		time.Sleep(1 * time.Second)
	}
}
