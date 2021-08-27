package main

import (
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	config *Config
}

func (server *Server) Run(handler http.Handler) error {
	addr := net.JoinHostPort("", server.config.ServerPort)
	if err := http.ListenAndServe(addr, handler); err != nil {
		return err
	}

	fmt.Printf("Server running on [%s]\n\n\n", addr)
	return nil
}

// NewServer returns a new Server instance
func NewServer(config *Config) *Server {
	return &Server{config}
}
