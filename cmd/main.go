package main

import (
	"fmt"
	"log"
	"net/http"

	"eliasschneider.com/cd-dc/cmd/config"
	"eliasschneider.com/cd-dc/cmd/docker"
	"eliasschneider.com/cd-dc/cmd/web"
)

func main() {
	addr := ":1411"
	http.HandleFunc("/upgrade/", handler)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	serviceName := r.URL.Path[len("/upgrade/"):]
	ctx := web.NewRequestContext(serviceName, r)
	request := ctx.Value("RequestContext").(web.RequestContext)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		web.HttpResponse(http.StatusMethodNotAllowed, "Method not allowed", &w, ctx)
		return
	}

	if r.Header.Get("X-Api-Key") != config.GetApiKey() {
		web.HttpResponse(http.StatusUnauthorized, "Unauthorized", &w, ctx)
		return
	}

	err := docker.UpgradeDockerComposeStack(ctx)
	if err == nil {
		web.HttpResponse(http.StatusOK, fmt.Sprintf("Service %s upgraded successfully", serviceName), &w, ctx)
	} else {
		web.HttpResponse(http.StatusInternalServerError, "An error occured while upgrading your service. Check the server logs.", &w, ctx)
		request.Logger.Printf("Error upgrading service %s", err.Error())
		return
	}
}
