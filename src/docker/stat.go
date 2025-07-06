package docker

import "time"

type Stat struct {
	Name			string		`json:"name"`
	Id				string		`json:"id"`
	ReadTime		time.Time	`json:"read"`
	CpuStat			CpuStat		`json:"cpu_stats"`
	PreCpuStat		CpuStat		`json:"precpu_stats"`
	CpuPercentage	float64		`json:"omitempty"`
}

type CpuStat struct {
	CpuUsage		CpuUsage	`json:"cpu_usage"`
	SysCpuUsage		uint64		`json:"system_cpu_usage"`
	OnlineCpu		uint8		`json:"online_cpus"`
} 

type CpuUsage struct {
	TotalUsage		 uint64		`json:"total_usage"`;
}

func (ds *Stat) CalculateCPUPercentageStream() {
	var containerCpu, systemCpu uint64;

	containerCpu = ds.CpuStat.CpuUsage.TotalUsage - ds.PreCpuStat.CpuUsage.TotalUsage;
	systemCpu = ds.CpuStat.SysCpuUsage - ds.PreCpuStat.SysCpuUsage;


	ds.CpuPercentage = (float64(containerCpu) / float64(systemCpu)) * 100
}