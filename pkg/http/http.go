package http

import (
	"gotify/config"
	"net/http"
	"time"
)

func NewHttp(cfg *config.Config) *http.Server {

	return &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      nil, // nil uses the default mux (http.DefaultServeMux)
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
