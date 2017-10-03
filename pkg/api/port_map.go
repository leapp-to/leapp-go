package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type PortMapParams struct {
	SourceHost       string              `json:"source_host"`
	TargetHost       string              `json:"target_host"`
	TcpPorts         TCPPortsUserMapping `json:"tcp_ports"`
	ExcludedTcpPorts ExcludedTCPPorts    `json:"excluded_tcp_ports"`
	DefaultPortMap   bool                `json:"default_port_map"`
}

func PortMap(request *http.Request) (interface{}, error) {
	var portMapParams PortMapParams

	if err := json.NewDecoder(request.Body).Decode(&portMapParams); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"source_host":            ObjValue{portMapParams.SourceHost},
		"target_host":            ObjValue{portMapParams.TargetHost},
		"use_default_port_map":   ObjValue{portMapParams.DefaultPortMap},
		"tcp_ports_user_mapping": portMapParams.TcpPorts,
		"excluded_tcp_ports":     portMapParams.ExcludedTcpPorts}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	exec := executor.New("port-mapping", string(json_data))
	result := exec.Execute()

	log.Println(result.Stderr)

	var out interface{}
	json.Unmarshal([]byte(result.Stdout), &out)
	return out, err
}
