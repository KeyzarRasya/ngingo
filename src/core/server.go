package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/KeyzarRasya/ngingo/src/cpu"
	"github.com/KeyzarRasya/ngingo/src/docker"
	"github.com/KeyzarRasya/ngingo/src/model"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Server struct {
	NgingoConfiguration NgingoConfiguration;
	DockerClient		*client.Client;
}

type PortCPU struct {
	Port 	uint16;
	CPU	float64;
}


func (s *Server) Run(config model.Configuration) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World | %s", r.URL)
	})

	pcpu, err := s.HighestCPU()

	if err != nil {
		return
	}

	fmt.Printf("Highest is PORT %d with Usage of : %f \n", pcpu.Port, pcpu.CPU)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)

	fmt.Printf("Running at Port %s", config.Port)

	for {}
}

func (s *Server) HighestCPU() (*PortCPU, error) {
	start := time.Now()
	var sortedPortCpus []PortCPU;
	portCpus, err := s.portAndCPU()

	if err != nil {
		return nil, err
	}

	for port, cpu := range portCpus {
		var portCpu PortCPU = PortCPU{
			Port: port,
			CPU: cpu,
		}
		sortedPortCpus = append(sortedPortCpus, portCpu)	
	}

	sortPortCPU(&sortedPortCpus)

	duration := time.Since(start)
	fmt.Printf("Reading took %s\n", duration)

	return &sortedPortCpus[0], nil

}

func sortPortCPU(portCPU *[]PortCPU) error {
	sort.Slice(*portCPU, func(i, j int) bool {
		return (*portCPU)[i].CPU > (*portCPU)[j].CPU 
	})
	return nil
}


func (s *Server) portAndCPU() (map[uint16]float64, error) {
	var portCpu map[uint16]float64 = make(map[uint16]float64)
	before, err := s.portAndStat()

	if err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Millisecond)

	after, err := s.portAndStat()

	if err != nil {
		return nil, err
	}

	if len(before) != len(after) {
		return nil, errors.New("There is might be one Port Down")
	}

	for port, stat := range before {
		portCpu[port] = docker.ReadCPUUsage(stat, after[port])
	}

	fmt.Printf("%f\n",portCpu[3000])
	fmt.Printf("%f\n",portCpu[3001])
	fmt.Printf("%f\n",portCpu[3002])

	return portCpu, nil

	
}

func (s *Server) portAndStat() (map[uint16]*cpu.Stat, error) {
	containers, err := s.DockerClient.ContainerList(context.Background(), container.ListOptions{})
	var portStat map[uint16]*cpu.Stat = make(map[uint16]*cpu.Stat);


	portCh := make(chan uint16, len(containers))
	statCh := make(chan *cpu.Stat, len(containers))


	if err != nil {
		return nil, err;
	}

	for _, ctr := range containers {
		go docker.ReadStatCh(s.DockerClient, ctr, portCh, statCh)
	}

	for i := 0; i < len(containers); i++ {

		stat := <-statCh
		port := <- portCh

		portStat[port] = stat;

	}

	return portStat, nil

}