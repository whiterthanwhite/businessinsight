package middleware

import (
	"log"
	"net/http"
)

type ServerLogger struct {
	Handler http.Handler
}

func (s *ServerLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request URI: %s\n", req.RequestURI)
	log.Printf("Method: %s\n", req.Method)

	s.Handler.ServeHTTP(w, req)
}
