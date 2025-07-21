package core

import (
	"fmt"
	"io"
	"math"

	"net/http"
	"net/url"

	"github.com/KeyzarRasya/ngingo/internal/balancer"
	"github.com/KeyzarRasya/ngingo/internal/docker"
	"github.com/KeyzarRasya/ngingo/internal/files"
	"github.com/KeyzarRasya/ngingo/internal/model"
)

type Server interface {
	Run()
}

type WebServer struct {
	NgingoConfiguration NgingoConfiguration
	Service             *docker.DockerService
	HttpClient          *http.Client
	Balancer            balancer.Balancer
	DataFiles           files.DataWriteRead
	Config              model.Configuration
}

func NewWebServer(
	ngingoConfiguration NgingoConfiguration,
	service *docker.DockerService,
	httpClient *http.Client,
	balancer balancer.Balancer,
	dataFiles files.DataWriteRead,
	config model.Configuration,
) WebServer {
	return WebServer{
		NgingoConfiguration: ngingoConfiguration,
		Service: service,
		HttpClient: httpClient,
		Balancer: balancer,
		DataFiles: dataFiles,
		Config: config,
	}
}

func (s *WebServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pcpu, err := s.Balancer.LowestUsage()
		if err != nil {
			return
		}

		before, err := s.Balancer.Usages()
		if err != nil {
			fmt.Println(err.Error())
			return
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
			return
		}

		diff := math.Abs(after[pcpu.Port] - before[pcpu.Port])
		if err := s.DataFiles.FormatAndWrite(pcpu.Port, pcpu.Endpoint, diff); err != nil {
			fmt.Println("Failed to write CPU Usage")
			return
		}

		s.SendResponse(w, string(b))
	})

	fmt.Printf("Running at Port %s", s.Config.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", s.Config.Port), nil)

	for {}
}

func (s *WebServer) SendResponse(w http.ResponseWriter, message string) {
	fmt.Fprint(w, message)
}

func (s *WebServer) processRequest(r *http.Request, portCpu *balancer.EndpointCPUStat) ([]byte, error) {
	url, err := url.Parse(r.URL.String())

	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%d%s", portCpu.Port, url.Path)

	req, err := http.NewRequest(r.Method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header = r.Header.Clone()
	res, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil

}
