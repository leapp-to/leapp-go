package api

import (
	"encoding/json"
	"net/http"
)

type destroyContainerParams struct {
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
	ContainerName string `json:"container_name"`
	Debug         bool   `json:"debug"`
}

func destroyContainer(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params destroyContainerParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	d := map[string]interface{}{
		"container_name":   ObjValue{params.ContainerName},
		"target_host":      ObjValue{params.TargetHost},
		"target_user_name": ObjValue{params.TargetUser},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "could not build actor's input")
	}

	id := actorRunnerRegistry.Create("destroy-container", string(actorInput))

	s := actorRunnerRegistry.GetStatus(id, true)
	if err := checkSyncTaskStatus(s); err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}
