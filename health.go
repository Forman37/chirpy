package main

import (
	"net/http"
)

type healthHandler struct {}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte("OK")
	w.Write(body)
}
