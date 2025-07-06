package balancer

type Balancer interface {
	LowestUsage() 	(*EndpointCPUStat, error);
	Usages()		(map[uint16]float64, error)
	
}

type VarStat interface {
	SetPortVarStat(uint16, string, float64)
	GetPortVarStat()	(uint16, string, float64)
	Clone()				VarStat
}
