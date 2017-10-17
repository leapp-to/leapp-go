package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type portMapParams struct {
	SourceHost          string              `json:"source_host"`
	SourceUser          string              `json:"source_user"`
	TargetHost          string              `json:"target_host"`
	TargetUser          string              `json:"target_user"`
	TCPPortsUserMapping TCPPortsUserMapping `json:"tcp_ports"`
	ExcludedTCPPorts    ExcludedTCPPorts    `json:"excluded_tcp_ports"`
	UseDefaultPortMap   bool                `json:"default_port_map"`
	Debug               bool                `json:"debug"`
}

func portMapHandler(request *http.Request) (*executor.Command, error) {
	var params portMapParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	d := map[string]interface{}{
		"source_host":            ObjValue{params.SourceHost},
		"target_host":            ObjValue{params.TargetHost},
		"use_default_port_map":   ObjValue{params.UseDefaultPortMap},
		"tcp_ports_user_mapping": params.TCPPortsUserMapping,
		"excluded_tcp_ports":     params.ExcludedTCPPorts,
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	c := executor.New("port-mapping", string(actorInput))
	return c, nil
}
