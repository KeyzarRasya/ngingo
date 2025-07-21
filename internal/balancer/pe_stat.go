package balancer

type EndpointCPUStat struct {
	Port 		uint16;
	CPU			float64;
	Endpoint	string;
}

func NewEndpointCPUStat() EndpointCPUStat {
	return EndpointCPUStat{}
}

func (pcpu *EndpointCPUStat) GetPortVarStat() (uint16, string, float64) {
	return pcpu.Port, pcpu.Endpoint, pcpu.CPU
}

func (pcpu *EndpointCPUStat) SetPortVarStat(port uint16, endpoint string, usage float64) {
	pcpu.Port = port;
	pcpu.Endpoint = endpoint
	pcpu.CPU = usage
}

func (pcpu *EndpointCPUStat) Clone() VarStat {
	copy := *pcpu
	return &copy
}

