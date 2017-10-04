package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type portMapParams struct {
	SourceHost       string              `json:"source_host"`
	TargetHost       string              `json:"target_host"`
	TCPPorts         TCPPortsUserMapping `json:"tcp_ports"`
	ExcludedTCPPorts ExcludedTCPPorts    `json:"excluded_tcp_ports"`
	DefaultPortMap   bool                `json:"default_port_map"`
}

func portMapHandler(request *http.Request) (interface{}, error) {
	var params portMapParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	d := map[string]interface{}{
		"source_host":            ObjValue{params.SourceHost},
		"target_host":            ObjValue{params.TargetHost},
		"use_default_port_map":   ObjValue{params.DefaultPortMap},
		"tcp_ports_user_mapping": params.TCPPorts,
		"excluded_tcp_ports":     params.ExcludedTCPPorts,
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	c := executor.New("port-mapping", string(actorInput))
	r := c.Execute()

	log.Println(r.Stderr)

	var out interface{}
	json.Unmarshal([]byte(r.Stdout), &out)
	return out, err
}
