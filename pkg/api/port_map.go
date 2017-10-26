package api

import (
	"encoding/json"
	"net/http"
)

type portMapParams struct {
	SourceHost       string              `json:"source_host"`
	TargetHost       string              `json:"target_host"`
	TCPPorts         TCPPortsUserMapping `json:"tcp_ports"`
	ExcludedTCPPorts ExcludedTCPPorts    `json:"excluded_tcp_ports"`
	DefaultPortMap   bool                `json:"default_port_map"`
}

func portMap(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params portMapParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, NewApiError(err, errBadInput, "could not decode data sent by client")
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
		return nil, http.StatusInternalServerError, NewApiError(err, errInternal, "could not build actor's input")
	}

	id := actorRunnerRegistry.Create("port-mapping", string(actorInput))

	s := actorRunnerRegistry.GetStatus(id, true)
	hs, err := checkTaskStatus(s)
	if err != nil {
		return nil, hs, NewApiError(err, errInternal, "could not build actor's input")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}
