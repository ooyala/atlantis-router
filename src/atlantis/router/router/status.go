package router

import (
	"atlantis/router/logger"
	"fmt"
	"net/http"
	"time"
)

type StatusServer struct {
	router  *Router
	backoff time.Duration
}

func NewStatusServer(r *Router) *StatusServer {
	return &StatusServer{
		router:  r,
		backoff: 1 * time.Second,
	}
}

func (s *StatusServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.config.Printer(w, r)
}

func (s *StatusServer) Run(port uint16, tout time.Duration) {
	server := http.Server{
		Handler:      s,
		Addr:         fmt.Sprintf("0.0.0.0:%u", port),
		ReadTimeout:  tout,
		WriteTimeout: tout,
	}
	logger.Errorf("[status server] %s", server.ListenAndServe())

	// try re-launching the status server
	s.backoff = s.backoff * 2
	go s.Run(port, tout)
}
