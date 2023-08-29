package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"eliasschneider.com/cd-dc/cmd/config"
	"eliasschneider.com/cd-dc/cmd/docker"
)

func main() {
	addr := ":1411"
	http.HandleFunc("/upgrade/", handler)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------- " + r.URL.Path)
	defer func() {
		fmt.Println("--------------------")
	}()

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		httpResponse(http.StatusMethodNotAllowed, "Method not allowed", &w)
		return
	}

	if r.Header.Get("X-Api-Key") != config.GetApiKey() {
		httpResponse(http.StatusUnauthorized, "Unauthorized", &w)
		return
	}

	serviceName := r.URL.Path[len("/upgrade/"):]

	err := docker.UpdateDockerComposeStack(serviceName)
	if err == nil {
		httpResponse(http.StatusOK, fmt.Sprintf("Service %s upgraded successfully", serviceName), &w)
	} else {
		httpResponse(http.StatusInternalServerError, "An error occured while upgrading your service. Check the server logs.", &w)
		log.Printf("Error updating %s: %s", serviceName, err.Error())
		return
	}
}

type HttpResponse struct {
	State   string `json:"state"`
	Message string `json:"error"`
}

func httpResponse(code int, message string, w *http.ResponseWriter) {
	writer := *w
	var state string
	if code == http.StatusOK {
		state = "success"
	} else {
		state = "failed"
	}

	response := HttpResponse{
		State:   state,
		Message: message,
	}

	responseJSON, _ := json.Marshal(response)
	writer.WriteHeader(code)
	writer.Write(responseJSON)

}
