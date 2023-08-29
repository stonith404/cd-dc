package main

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"slices"
	"strings"
)

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

	services := GetServices()
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

	log.Printf("Finished updating %s", serviceName)

	return err
}


func PullImages(service Service, name string) error {
	containersFormatted := strings.Join(service.Containers, " ")
	err := runCommand("compose", "-f", service.Path, "pull", containersFormatted)
	return err
}

func RestartContainers(service Service, name string) error {
	containersFormatted := strings.Join(service.Containers, " ")
	err := runCommand("compose", "-f", service.Path, "up", containersFormatted, "-d")
	return err
}


func runCommand(command ...string) error {

	// remove empty strings from command
	command = slices.DeleteFunc(command, func(s string) bool {
		return s == ""
	})

	cmd := exec.Command("docker", command...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}
