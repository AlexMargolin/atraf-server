package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseData struct {
	Data interface{} `json:"data"`
}

// Response is a generic and unified way to return a response to the client.
// Should be used when the response is considered "successful".
// When we would like to respond with an error, it's better to just use http.Error with the
// appropriate error http.status code.
func Response(w http.ResponseWriter, code int, data ...interface{}) {
	response := ResponseData{data}

	if len(data) > 0 {
		response.Data = data[0]
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	// Set Headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Http status code
	w.WriteHeader(code)

	// When write fails, we most likely won't be able to respond to the client.
	if _, err = w.Write(responseJson); err != nil {
		log.Println(err)
		return
	}
}
