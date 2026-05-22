package httpserver

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	wsUpd      *websocket.Upgrader
}

func NewHTTPServer(host, port string, timeout, idletimeout time.Duration) *Server {
	srv := &http.Server{
		Addr:           host + ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		IdleTimeout:    idletimeout,
	}
	ws := websocket.Upgrader{EnableCompression: true}
	return &Server{
		httpServer: srv,
		wsUpd:      &ws,
	}
}
func (s *Server) Upgrader() *websocket.Upgrader {
	return s.wsUpd
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer.Handler = handler
	return s.httpServer.ListenAndServe()
}

func (s *Server) Gracefull(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
