package app

import (
	"net"
	"net/http"
)

type ServerConfig struct {
	Host string
	Port string
}

func RunServer(c *ServerConfig, h http.Handler) error {
	addr := net.JoinHostPort(c.Host, c.Port)
	if err := http.ListenAndServe(addr, h); err != nil {
		return err
	}
	return nil
}
