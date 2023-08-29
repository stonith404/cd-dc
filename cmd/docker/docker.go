package docker

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"slices"
	"strings"
	"time"

	"eliasschneider.com/cd-dc/cmd/config"
	"eliasschneider.com/cd-dc/cmd/web"
)

var runningJobs = []string{}

func UpgradeDockerComposeStack(ctx context.Context) error {
	request := ctx.Value("RequestContext").(*web.RequestContext)

	alreadyUpdating := slices.Contains(runningJobs, request.ServiceName)
	defer func() {
		if !alreadyUpdating {
			runningJobs = slices.DeleteFunc(runningJobs, func(s string) bool {
				return s == request.ServiceName
			})
		}
	}()
	var err error

	services := config.GetServices()
	service, exists := services[request.ServiceName]
	if !exists {
		return errors.New("service not found")
	}

	// If the service is already updating, wait 5 seconds and try again
	if alreadyUpdating {
		request.Logger.Print("Already updating service, waiting")
		time.Sleep(5 * time.Second)
		return UpgradeDockerComposeStack(ctx)
	}

	runningJobs = append(runningJobs, request.ServiceName)
	if service.Path == "" {
		return errors.New("service path not found")
	}

	request.Logger.Print("Pulling images")
	if err := PullImages(service, request.ServiceName); err != nil {
		return err
	}

	request.Logger.Print("Recreating containers")
	if err := RestartContainers(service, request.ServiceName); err != nil {
		return err
	}

	request.Logger.Print("Pruning old images")
	err = PruneOldImages(ctx, service)
	if err != nil {
		return err
	}

	request.Logger.Print("Finished updating")

	return err
}

func PullImages(service config.Service, name string) error {
	containersFormatted := strings.Join(service.Containers, " ")
	_, err := runCommand("docker", "compose", "-f", service.Path, "pull", containersFormatted)
	return err
}

func RestartContainers(service config.Service, name string) error {
	containersFormatted := strings.Join(service.Containers, " ")
	_, err := runCommand("docker", "compose", "-f", service.Path, "up", containersFormatted, "-d")
	return err
}

func runCommand(name string, command ...string) (string, error) {
	// remove empty strings from command
	command = slices.DeleteFunc(command, func(s string) bool {
		return s == ""
	})

	cmd := exec.Command(name, command...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}
