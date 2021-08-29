package app

import (
	"net"
	"net/http"
)

type ServerConfig struct {
	Host    string
	Port    string
	Handler http.Handler
}

func Run(c *ServerConfig) error {
	addr := net.JoinHostPort(c.Host, c.Port)
	if err := http.ListenAndServe(addr, c.Handler); err != nil {
		return err
	}
	return nil
}
