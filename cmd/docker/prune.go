package docker

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"eliasschneider.com/cd-dc/cmd/config"
)

// Remove dangling images for a given image name.
// We keep the newest image for rollbacking.
func PruneOldImages(service config.Service) error {
	imageNames, err := getDockerImageNamesOfService(service)
	if err != nil {
		return err
	}

	for _, imageName := range imageNames {
		imageIdsRaw, err := runCommand("docker", "images", "-f", "dangling=true", "-q", imageName)
		if err != nil {
			return fmt.Errorf("error getting dangling images for %s: %s", imageName, err.Error())
		}

		var imageIds []string
		if imageIdsRaw != "" {
			imageIds = strings.Split(strings.TrimSpace(imageIdsRaw), "\n")
		}

		if len(imageIds) == 0 {
			log.Printf("No dangling images found for %s", imageName)
			return nil
		}
		// Remove newest image
		imageIds = slices.Delete(imageIds, 0, config.GetNumberOfImagesToKeep()-1)

		for _, imageId := range imageIds {
			log.Printf("Removing dangling image %s", imageId)
			_, err := runCommand("docker", "rmi", imageId)
			if err != nil {
				return fmt.Errorf("error removing image %s: %s", imageId, err.Error())
			}
		}
	}

	return nil
}

func getDockerImageNamesOfService(service config.Service) ([]string, error) {

	imageNames := []string{}

	if service.Containers == nil {
		imageNamesRaw, err := runCommand("/bin/sh", "-c", "grep 'image:' "+service.Path+" | awk '{print $2}'")
		if err != nil {
			return nil, err
		}
		imageNames = strings.Split(strings.TrimSpace(imageNamesRaw), "\n")
		for i, imageName := range imageNames {
			imageNames[i] = strings.Split(imageName, ":")[0]
		}

	} else {

		for _, container := range service.Containers {
			imageName, err := runCommand("/bin/sh", "-c", "grep -m 1 -A 1 "+container+" "+service.Path+" | grep 'image:' | awk '{print $2}'")
			if err != nil {
				return nil, err
			}
			imageName = strings.Split(strings.TrimSpace(imageName), ":")[0]
			imageNames = append(imageNames, imageName)
		}
	}

	return imageNames, nil

}
