package cpu

import (
	"errors"
	"sort"
	"time"

	"github.com/KeyzarRasya/ngingo/src/balancer"
	"github.com/KeyzarRasya/ngingo/src/docker"
)

type CPUBalancer struct {
	Service docker.Service
}

func NewCPUBalancer(service docker.Service) CPUBalancer {
	return CPUBalancer{Service: service}
} 

func (cb *CPUBalancer) LowestUsage() (*balancer.EndpointCPUStat, error) {
	var sortedPortCpus []balancer.EndpointCPUStat;
	portCpus, err := cb.portAndCPU()

	if err != nil {
		return nil, err
	}

	for port, cpu := range portCpus {
		var portCpu balancer.EndpointCPUStat = balancer.EndpointCPUStat{
			Port: port,
			CPU: cpu,
		}
		sortedPortCpus = append(sortedPortCpus, portCpu)	
	}

	sortPortCPU(&sortedPortCpus)

	return &sortedPortCpus[0], nil

}

func (cb *CPUBalancer) Usages() (map[uint16]float64, error) {
	pcpu, err := cb.portAndCPU();

	return pcpu, err

}

func sortPortCPU(portCPU *[]balancer.EndpointCPUStat) error {
	sort.Slice(*portCPU, func(i, j int) bool {
		return (*portCPU)[i].CPU < (*portCPU)[j].CPU 
	})
	return nil
}

func (cb *CPUBalancer) portAndCPU() (map[uint16]float64, error) {
	var portCpu map[uint16]float64 = make(map[uint16]float64)
	before, err := cb.Service.ReadStat()

	if err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Millisecond)

	after, err := cb.Service.ReadStat()

	if err != nil {
		return nil, err
	}

	if len(before) != len(after) {
		return nil, errors.New("There is might be one Port Down")
	}

	
	for port, stat := range before {
		portCpu[port] = readCPUUsage(stat, after[port])
	}

	// fmt.Printf("%f\n",portCpu[3000])
	// fmt.Printf("%f\n",portCpu[3001])
	// fmt.Printf("%f\n",portCpu[3002])

	return portCpu, nil

}

func readCPUUsage(before, after *docker.Stat) float64 {
	cpuDelta := after.CpuStat.CpuUsage.TotalUsage - before.CpuStat.CpuUsage.TotalUsage
	sysDelta := after.CpuStat.SysCpuUsage - before.CpuStat.SysCpuUsage
	onlineCPU := after.CpuStat.OnlineCpu

	cpuPercent := (float64(cpuDelta) / float64(sysDelta)) * float64(onlineCPU) * 100

	return cpuPercent
}
