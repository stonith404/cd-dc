package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	State       string `json:"state"`
	Message     string `json:"error"`
	ServiceName string `json:"serviceName"`
	RequestID   string `json:"requestId"`
}

func HttpResponse(ctx context.Context, code int, message string, w *http.ResponseWriter) {
	request := ctx.Value("RequestContext").(RequestContext)

	writer := *w
	var state string
	if code == http.StatusOK {
		state = "success"
	} else {
		state = "failed"
	}

	response := HTTPResponse{
		State:       state,
		Message:     message,
		ServiceName: request.ServiceName,
		RequestID:   request.RequestID,
	}

	responseJSON, _ := json.Marshal(response)
	writer.WriteHeader(code)
	writer.Write(responseJSON)

}
