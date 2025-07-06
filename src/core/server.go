package core

import (
	"fmt"
	"io"

	"net/http"
	"net/url"

	"github.com/KeyzarRasya/ngingo/src/balancer"
	"github.com/KeyzarRasya/ngingo/src/docker"
	"github.com/KeyzarRasya/ngingo/src/files"
	"github.com/KeyzarRasya/ngingo/src/model"
)

type Server struct {
	NgingoConfiguration NgingoConfiguration;
	Service 			*docker.DockerService
	HttpClient			*http.Client;
	Balancer			balancer.Balancer;
	DataFiles			files.DataWriteRead
}

func (s *Server) Run(config model.Configuration) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var records [][]string;
		pcpu, err := s.Balancer.LowestUsage()
		if err != nil {
			return;
		}

		before, err := s.Balancer.Usages();
		if err != nil {
			fmt.Println(err.Error())
			return;
		}
		
		b, err := s.processRequest(r, pcpu)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		pcpu.Endpoint = r.URL.Path

		after, err := s.Balancer.Usages()
		if err != nil {
			fmt.Println(err.Error())
			return;
		}

		diff := (after[pcpu.Port] - before[pcpu.Port])

		record := []string{fmt.Sprintf("%d", pcpu.Port), pcpu.Endpoint, fmt.Sprintf("%f", diff)}
		records = append(records, record)

		s.DataFiles.Write(records)

		fmt.Printf("===\n[PORT : %d]\n Before Usage : %f\n After Usage :  %f\n Diff : %f\n===\n", pcpu.Port, before[pcpu.Port], after[pcpu.Port], diff)

	

		fmt.Fprint(w, string(b))

	})

	fmt.Printf("Running at Port %s", config.Port)

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)


	for {}
}


func (s *Server) processRequest(r *http.Request, portCpu *balancer.EndpointCPUStat) ([]byte, error) {
	url, err := url.Parse(r.URL.String())

	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%d%s", portCpu.Port, url.Path)

	req, err := http.NewRequest(r.Method, endpoint, nil);

	if err != nil {
		return nil, err
	}

	req.Header = r.Header.Clone()

	res, err := s.HttpClient.Do(req)

	if err != nil {
		return nil, err;
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return b, nil;
	
}