package middleware

import "net/http"

type ReactHelper struct {
	Handler http.Handler
}

func (rh *ReactHelper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		return
	}

	rh.Handler.ServeHTTP(w, req)
}
