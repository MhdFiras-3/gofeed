package handlers

import (
	"encoding/json"
	"net/http"

	"log"
)

func respWithError(w http.ResponseWriter, errStatus int, errMsg string) {
	type errorMsg struct {
		Error string `json:"error"`
	}

	data, err := json.Marshal(errorMsg{Error: errMsg})
	if err != nil {
		log.Printf("error encoding data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errStatus)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("failed to respond: %v", err)
	}
}

func respWithJson(w http.ResponseWriter, respStatus int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error encoding data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respStatus)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("failed to respond: %v", err)
	}
}
