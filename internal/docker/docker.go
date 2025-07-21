package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Service interface {
	AddContainer()	error;
	ReadStat()	(map[uint16]*Stat, error);
}

type DockerService struct {
	dockerInstance		*client.Client
}

func NewDockerService() (*DockerService, error) {
	client, err :=  client.NewClientWithOpts(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithVersion("1.50"),
	)

	if err != nil {
		return nil, err
	}

	return &DockerService{dockerInstance: client}, nil
}

func (ds *DockerService) AddContainer() error {
	return nil
}

func (ds *DockerService) ReadStat() (map[uint16]*Stat, error) {
	containers, err := ds.dockerInstance.ContainerList(context.Background(), container.ListOptions{})

	var portStat map[uint16]*Stat = make(map[uint16]*Stat);

	portChan := make(chan uint16, len(containers))
	statChan := make(chan *Stat, len(containers)) 

	if err != nil {
		return nil, err;
	}

	for _, ctr := range containers {
		go ds.readStatCh(ctr, portChan, statChan)
	} 

	for i := 0; i < len(containers); i++ {

		stat := <-statChan
		port := <- portChan

		portStat[port] = stat;

	}

	return portStat, nil
}

func (ds *DockerService) readStatCh(ctr container.Summary, port chan<- uint16, cpuStat chan<- *Stat) error {
	resp, err := ds.dockerInstance.ContainerStatsOneShot(context.Background(), ctr.ID)

	if err != nil {
		fmt.Printf("Failed to read container stats : %s", err.Error());
		port <- 0;
		cpuStat <- nil
		return err;
	}

	defer resp.Body.Close()

	var stat Stat;
	var memstat runtime.MemStats;

	info, err := io.ReadAll(resp.Body)
	if err != nil {
		return err;
	}

	runtime.ReadMemStats(&memstat);
	fmt.Println(string(info))

	err = json.Unmarshal(info, &stat);
	if err != nil {
		fmt.Println(err.Error())
		port <- 0;
		cpuStat <- nil
		return err;
	}

	port <- ctr.Ports[1].PublicPort
	cpuStat <- &stat
	return nil
}


