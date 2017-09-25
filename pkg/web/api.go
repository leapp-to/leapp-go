package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

type MigrateParams struct {
	SourceHost       string            `json:"source_host,omitempty"`
	TargetHost       string            `json:"target_target,omitempty"`
	ContainerName    string            `json:"container_name,omitempty"`
	SourceUser       string            `json:"source_user,omitempty"`
	TargetUser       string            `json:"target_user,omitempty"`
	ExcludePaths     []string          `json:"excluded_paths,omitempty"`
	TcpPorts         map[uint16]uint16 `json:"tcp_ports,omitempty"`
	ExcludedTcpPorts []uint16          `json:"excluded_tcp_ports,omitempty"`
	ForceCreate      bool              `json:"force_create,omitempty"`
	DisableStart     bool              `json:"disable_start,omitempty"`
	Debug            bool              `json:"debug,omitempty"`
}

func MigrateMachine(request *http.Request) (interface{}, error) {
	var (
		migrateParams MigrateParams
		data          interface{}
	)

	if err := json.NewDecoder(request.Body).Decode(&migrateParams); err != nil {
		return nil, err
	}

	// place for executor
	// Executor.execute // or something
	output_mock := []byte("{\"foo\": \"bar\"}")
	// return result to client

	if json.Unmarshal(output_mock, &data) != nil {
		return nil, errors.New("Invalid json output from executor")
	}

	result := data.(map[string]interface{})
	return result, nil

}
