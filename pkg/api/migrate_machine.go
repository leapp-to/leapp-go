package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

// migrateParams represents the data sent by the client.
type migrateParams struct {
	SourceHost          string              `json:"source_host"`
	SourceUser          string              `json:"source_user"`
	TargetHost          string              `json:"target_host"`
	TargetUser          string              `json:"target_user"`
	TCPPortsUserMapping TCPPortsUserMapping `json:"tcp_ports"`
	ExcludedTCPPorts    ExcludedTCPPorts    `json:"excluded_tcp_ports"`
	UseDefaultPortMap   bool                `json:"default_port_map"`
	ContainerName       string              `json:"container_name"`
	StartContainer      bool                `json:"start_container"`
	ForceCreate         bool                `json:"force_create"`
	ExcludePaths        []string            `json:"excluded_paths"`
	Debug               bool                `json:"debug"`
}

// buildActorInput translates the data sent by the client into data that the actor can interpret.
func buildActorInput(p *migrateParams) (string, error) {
	data := map[string]interface{}{
		"start_container":        ObjValue{p.StartContainer},
		"container_name":         ObjValue{p.ContainerName},
		"force_create":           ObjValue{p.ForceCreate},
		"source_host":            ObjValue{p.SourceHost},
		"source_user_name":       ObjValue{p.SourceUser},
		"target_host":            ObjValue{p.TargetHost},
		"target_user_name":       ObjValue{p.TargetUser},
		"excluded_paths":         ObjValue{p.ExcludePaths},
		"use_default_port_map":   ObjValue{p.UseDefaultPortMap},
		"tcp_ports_user_mapping": p.TCPPortsUserMapping,
		"excluded_tcp_ports":     p.ExcludedTCPPorts,
	}

	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

// migrateMachineHandler handles the /migrate-machine endpoint.
func migrateMachineHandler(request *http.Request) (*executor.Command, error) {
	var params migrateParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	// Translate data sent by client into data that actor can read
	actorInput, err := buildActorInput(&params)
	if err != nil {
		return nil, err
	}

	// Creates an executor.Command that calls the correct actor passing the data to its stdin
	c := executor.New("migrate-machine", actorInput)

	return c, nil
}
