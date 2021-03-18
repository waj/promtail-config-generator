package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type promtailConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func main() {
	environment, found := os.LookupEnv("ENV")
	if !found {
		log.Fatalln("ENV environment variable missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := client.NewEnvClient()

	if err != nil {
		log.Fatal(err)
	}

	opts := types.ContainerListOptions{}
	containers, err := client.ContainerList(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	configs := []promtailConfig{}

	for _, container := range containers {
		containerInspect, err := client.ContainerInspect(ctx, container.ID)
		if err != nil {
			log.Fatal(err)
		}

		rancherStackService := containerInspect.Config.Labels["io.rancher.stack_service.name"]
		if rancherStackService != "" {
			stackAndService := strings.Split(rancherStackService, "/")

			config := promtailConfig{
				Targets: []string{"localhost"},
				Labels: map[string]string{
					"__path__": containerInspect.LogPath,
					"stack":    stackAndService[0],
					"service":  rancherStackService,
					"env":      environment,
				},
			}

			configs = append(configs, config)
		}
	}

	configJson, err := json.Marshal(configs)
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Create("/etc/promtail-rancher/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	configFile.Write(configJson)
}
