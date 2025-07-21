package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/KeyzarRasya/ngingo/internal/balancer"
	"github.com/KeyzarRasya/ngingo/internal/balancer/cpu"
	"github.com/KeyzarRasya/ngingo/internal/core"
	"github.com/KeyzarRasya/ngingo/internal/docker"
	"github.com/KeyzarRasya/ngingo/internal/files"
	"github.com/KeyzarRasya/ngingo/internal/model"
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

	if err := parseYAML(model.PATH ,&config); err != nil {
		fmt.Println(err)
		return;
	}

	if err := parseYAML(config.ConfigPath, &ngingo); err != nil {
		fmt.Println(err)
		return;
	}

	ds, err := docker.NewDockerService()
	
	if err != nil {
		fmt.Printf("Failed to create docker client")
		return;
	}
	
	cpuBalancer := cpu.NewCPUBalancer(ds)
	endpointStat := balancer.NewEndpointCPUStat()
	dataCpu := files.NewDataCPU(config.FileCPU, &endpointStat)

	server := core.NewWebServer(
		ngingo,
		ds,
		&http.Client{},
		&cpuBalancer,
		&dataCpu,
		config,
	)

	server.Run()

}