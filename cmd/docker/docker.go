package docker

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"slices"
	"strings"

	"eliasschneider.com/cd-dc/cmd/config"
)

var runningJobs = []string{}

func UpdateDockerComposeStack(serviceName string) error {
	alreadyUpdating := slices.Contains(runningJobs, serviceName)
	defer func() {
		if !alreadyUpdating {
			runningJobs = slices.DeleteFunc(runningJobs, func(s string) bool {
				return s == serviceName
			})
		}
	}()
	var err error

	services := config.GetServices()
	service, exists := services[serviceName]
	if !exists {
		return errors.New("service not found")
	}

	if alreadyUpdating {
		log.Printf("Already updating %s, skipping", serviceName)
		return errors.New("already updating")
	}

	runningJobs = append(runningJobs, serviceName)
	if service.Path == "" {
		log.Default().Printf("path is empty for service %s", serviceName)
	}

	log.Printf("Pulling images for %s", serviceName)
	if err := PullImages(service, serviceName); err != nil {
		return err
	}

	log.Printf("Recreating containers for %s", serviceName)
	if err := RestartContainers(service, serviceName); err != nil {
		return err
	}

	log.Printf("Pruning old images for %s", serviceName)
	err = PruneOldImages(service)
	if err != nil {
		return err
	}

	log.Printf("Finished updating %s", serviceName)

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
