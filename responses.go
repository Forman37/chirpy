package main

import (
	"log"
	"encoding/json"
	"net/http"
)

type response struct {
	Error string `json:"error"`
	Valid bool `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error){
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	responseBody := response{
		Error: msg,
	}

	dat, err2 := json.Marshal(responseBody)
	if err2 != nil {
		log.Printf("Error marshalling JSON: %s\n", err2)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling OK response: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
