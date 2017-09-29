package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

// RUNNER is an actor execution tool
// TODO: using a wrapper script for testing only, but runner.py in snactor repo should be refactored into a standalone tool so it can be used by leapp-daemon
const RUNNER = "/home/fjb/src/snactor/runner_wrapper.sh"

// MigrateParams represents the parameters sent by the client
// TODO: there are more parameters to add here
type MigrateParams struct {
	SourceHost       string            `json:"source_host"`
	TargetHost       string            `json:"target_host"`
	ContainerName    string            `json:"container_name"`
	SourceUser       string            `json:"source_user"`
	TargetUser       string            `json:"target_user"`
	ExcludePaths     []string          `json:"excluded_paths"`
	TCPPorts         map[uint16]uint16 `json:"tcp_ports"`
	ExcludedTCPPorts []uint16          `json:"excluded_tcp_ports"`
	ForceCreate      bool              `json:"force_create"`
	DisableStart     bool              `json:"disable_start"`
	Debug            bool              `json:"debug"`
}

// buildActorInput translates the data sent by the client into data that the actor can interpret.
// TODO: this function is not ready yet; it should validate input data and set default values
func buildActorInput(p *MigrateParams) (string, error) {
	data := make(map[string]interface{})

	var sc bool
	if p.DisableStart == true {
		sc = false
	} else {
		sc = true
	}
	data["start_container"] = map[string]interface{}{"value": sc}
	data["container_name"] = map[string]interface{}{"value": p.ContainerName}
	data["force_create"] = map[string]interface{}{"value": p.ForceCreate}
	data["source_host"] = map[string]interface{}{"value": p.SourceHost}
	data["source_user_name"] = map[string]interface{}{"value": p.SourceUser}
	data["target_host"] = map[string]interface{}{"value": p.TargetHost}
	data["target_user_name"] = map[string]interface{}{"value": p.TargetUser}
	data["excluded_paths"] = map[string]interface{}{"value": p.ExcludePaths}
	data["excluded_tcp_ports"] = map[string]interface{}{"tcp": make(map[string]string)}
	data["tcp_ports_user_mapping"] = map[string]interface{}{"ports": []string{}}
	data["use_default_port_map"] = map[string]interface{}{"value": true}

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
	actorData, err := buildActorInput(&params)
	if err != nil {
		return nil, err
	}

	// Call the actor runner passing data to its stdin
	c := executor.Command{
		CmdLine: strings.Split(RUNNER, " "),
		Stdin:   actorData,
	}

	// TODO: this blocks until the migration is done; so we might execute this in a worker in the future
	r := c.Execute()

	log.Println(r.Stderr)

	return r.Stdout, nil
}
