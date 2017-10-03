package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

// MigrateParams represents the parameters sent by the client
type MigrateParams struct {
	StartContainer      bool                `json:"start_container"`
	ContainerName       string              `json:"container_name"`
	ForceCreate         bool                `json:"force_create"`
	SourceHost          string              `json:"source_host"`
	SourceUser          string              `json:"source_user"`
	TargetHost          string              `json:"target_host"`
	TargetUser          string              `json:"target_user"`
	ExcludePaths        []string            `json:"excluded_paths"`
	UseDefaultPortMap   bool                `json:"use_default_port_map"`
	TCPPortsUserMapping TCPPortsUserMapping `json:"tcp_ports_user_mapping"`
	ExcludedTCPPorts    ExcludedTCPPorts    `json:"excluded_tcp_ports"`
}

// buildActorInput translates the data sent by the client into data that the actor can interpret.
func buildActorInput(p *MigrateParams) (string, error) {
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

// MigrateMachine handles the /migrate-machine endpoint.
func MigrateMachine(request *http.Request) (interface{}, error) {
	var params MigrateParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	// Translate data sent by client into data that actor can read
	actorInput, err := buildActorInput(&params)
	if err != nil {
		return nil, err
	}

	// Call the actor runner passing data to its stdin
	c := executor.New("migrate-machine", actorInput)

	// TODO: this blocks until the migration is done; so we might execute this in a worker in the future
	r := c.Execute()

	log.Println(r.Stderr)

	return r.Stdout, nil
}
