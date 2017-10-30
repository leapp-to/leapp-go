package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// migrateParams represents the data sent by the client.
type migrateParams struct {
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

func migrateMachineStart(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var p migrateParams
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	actorInput, err := buildActorInput(&p)
	if err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "could not build actor's input")
	}

	id := actorRunnerRegistry.Create("migrate-machine", actorInput)

	return map[string]*ActorRunnerID{"migrate-id": id}, http.StatusOK, nil
}

func migrateMachineResult(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	id, err := parseID(req)
	if err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not parse ID")
	}

	s := actorRunnerRegistry.GetStatus(id, false)
	if s == nil {
		return nil, http.StatusNotFound, newAPIError(nil, errTaskNotFound, "task not found")
	}

	if s.Result == nil {
		return nil, http.StatusOK, newAPIError(nil, errTaskRunning, "task found, but there is no result yet")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}

func parseID(req *http.Request) (*ActorRunnerID, error) {
	i, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, errors.New("ID not found in request")
	}

	aid, err := uuid.Parse(i)
	if err != nil {
		return nil, err
	}
	return &ActorRunnerID{aid}, nil
}
