package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseData struct {
	Data interface{} `json:"data"`
}

func Success(w http.ResponseWriter, code int, data ...interface{}) {
	response := ResponseData{data}

	if len(data) > 0 {
		response.Data = data[0]
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)

	// When write fails, we most likely won't be able to respond to the client.
	if _, err = w.Write(responseJson); err != nil {
		log.Println(err)
		return
	}
}

func Error(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
}
