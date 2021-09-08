package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseData struct {
	Data interface{} `json:"data"`
}

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func Success(w http.ResponseWriter, code int, data ...interface{}) {
	SetHeaders(w)

	// exit early for no content responses
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	response := ResponseData{data}
	if len(data) > 0 {
		response.Data = data[0]
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// When write fails, we most likely won't be able to respond to the client.
	if _, err = w.Write(encoded); err != nil {
		log.Println(err)
		return
	}
}

func Error(w http.ResponseWriter, code int) {
	SetHeaders(w)

	w.WriteHeader(code)
}
