package balancer

import "github.com/KeyzarRasya/ngingo/src/docker"

type Balancer interface {
	LowestUsage() 							(*EndpointCPUStat, error);
	Usages()								(map[uint16]float64, error)
	ReadUsage(before, after *docker.Stat)	float64
}

type VarStat interface {
	SetPortVarStat(uint16, string, float64)
	GetPortVarStat()	(uint16, string, float64)
	Clone()				VarStat
}
