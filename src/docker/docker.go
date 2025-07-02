package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KeyzarRasya/ngingo/src/cpu"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Init() (*client.Client, error) {
	return client.NewClientWithOpts(
			client.WithHost("unix:///var/run/docker.sock"),
			client.WithVersion("1.50"),
	)
}

func ReadStat(client *client.Client, ctr container.Summary) (*cpu.Stat, error) {
	resp, err := client.ContainerStats(context.Background(), ctr.ID, true)

	if err != nil {
		fmt.Printf("Failed to read container stats : %s", err.Error());
		return nil, err;
	}

	defer resp.Body.Close()

	
	var stat cpu.Stat;

	decoder := json.NewDecoder(resp.Body) 
	err = decoder.Decode(&stat)
	
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//fmt.Printf("[%s] Name : %s  | CPU Percentage : %f | Read : %s\n", state, stat.Name, stat.CpuPercentage, stat.ReadTime.String())

	return &stat, nil;

}

func ReadStatCh(client *client.Client, ctr container.Summary, port chan<- uint16, cpuStat chan<- *cpu.Stat) error {
	resp, err := client.ContainerStatsOneShot(context.Background(), ctr.ID)

	if err != nil {
		fmt.Printf("Failed to read container stats : %s", err.Error());
		port <- 0;
		cpuStat <- nil
		return err;
	}

	defer resp.Body.Close()

	var stat cpu.Stat;

	decoder := json.NewDecoder(resp.Body) 
	err = decoder.Decode(&stat)
	
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

func ReadCPUUsagePerRequest(client *client.Client, ctr container.Summary) error {
	before,  err := ReadStat(client, ctr)

	if err != nil {
		fmt.Println(err.Error())
		return err;
	}

	for range 3000 {
		_, err = http.Get(fmt.Sprintf("http://localhost:%d/", ctr.Ports[1].PublicPort))
	
		if err != nil {
			fmt.Println(err.Error())
			return err;
		}
	}

	after, err := ReadStat(client, ctr);

	if err != nil {
		fmt.Println(err.Error())
		return err;
	}

	cpuPercent := ReadCPUUsage(before, after);

	fmt.Printf("[%d] CPU used by request: %.6f%%\n", ctr.Ports[1].PublicPort, cpuPercent)

	return nil
}

func ReadCPUUsage(before, after *cpu.Stat) float64 {
	cpuDelta := after.CpuStat.CpuUsage.TotalUsage - before.CpuStat.CpuUsage.TotalUsage
	sysDelta := after.CpuStat.SysCpuUsage - before.CpuStat.SysCpuUsage
	onlineCPU := after.CpuStat.OnlineCpu

	cpuPercent := (float64(cpuDelta) / float64(sysDelta)) * float64(onlineCPU) * 100

	return cpuPercent
}

