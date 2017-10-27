package api

import (
	"encoding/json"
	"net/http"
)

type portInspectParams struct {
	TargetHost  string `json:"target_host"`
	PortRange   string `json:"port_range"`
	ShallowScan bool   `json:"shallow_scan"`
}

func portInspect(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params portInspectParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	d := map[string]interface{}{
		"host": ObjValue{params.TargetHost},
		"scan_options": map[string]interface{}{
			"shallow_scan": params.ShallowScan,
			"port_range":   params.PortRange,
			"force_nmap":   !params.ShallowScan,
		},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "could not build actor's input")
	}

	id := actorRunnerRegistry.Create("portscan", string(actorInput))

	s := actorRunnerRegistry.GetStatus(id, true)
	hs, err := checkTaskStatus(s)
	if err != nil {
		return nil, hs, newAPIError(err, errInternal, "could not build actor's input")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}
