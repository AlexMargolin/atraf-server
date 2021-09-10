package rest

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
}

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func Success(w http.ResponseWriter, code int, data interface{}) {
	SetHeaders(w)

	if data == nil {
		w.WriteHeader(code)
		return
	}

	response := Response{data}
	encoded, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// When write fails, we most likely won't be able to respond to the client.
	if _, err = w.Write(encoded); err != nil {
		// TODO add error log
		return
	}
}

func Error(w http.ResponseWriter, code int) {
	SetHeaders(w)

	w.WriteHeader(code)
}
