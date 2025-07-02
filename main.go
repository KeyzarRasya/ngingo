package main

import (
	"fmt"
	"os"

	"github.com/KeyzarRasya/ngingo/src/core"
	"github.com/KeyzarRasya/ngingo/src/docker"
	"github.com/KeyzarRasya/ngingo/src/model"
	"gopkg.in/yaml.v3"
)

func parseYAML(path string, out interface{}) error {
	b, err := os.ReadFile(path);

	if err != nil {
		fmt.Println("Failed to read config.yaml file");
		return err;
	}

	if err := yaml.Unmarshal(b, out); err != nil {
		fmt.Println("Failed to parse config.yaml");
		return err;
	}

	return nil;
}

func main() {
	var config model.Configuration;
	var ngingo core.NgingoConfiguration;
	var server core.Server;

	if err := parseYAML(model.PATH ,&config); err != nil {
		fmt.Println(err)
		return;
	}

	if err := parseYAML(config.ConfigPath, &ngingo); err != nil {
		fmt.Println(err)
		return;
	}

	
	client, err := docker.Init()
	
	if err != nil {
		fmt.Printf("Failed to create docker client")
		return;
	}
	
	server = core.Server{
		NgingoConfiguration: ngingo,
		DockerClient: client,
	}
	

	// containers, err := client.ContainerList(context.Background(), container.ListOptions{});

	// for _, ctr := range containers {
	// 	go docker.ReadStat(client, ctr)
	// }

	server.Run(config)

}