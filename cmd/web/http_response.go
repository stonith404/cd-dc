package web

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	State   string `json:"state"`
	Message string `json:"error"`
}

func HttpResponse(code int, message string, w *http.ResponseWriter) {
	writer := *w
	var state string
	if code == http.StatusOK {
		state = "success"
	} else {
		state = "failed"
	}

	response := HTTPResponse{
		State:   state,
		Message: message,
	}

	responseJSON, _ := json.Marshal(response)
	writer.WriteHeader(code)
	writer.Write(responseJSON)

}
