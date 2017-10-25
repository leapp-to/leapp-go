package api

import (
	"encoding/json"
	"net/http"
)

type destroyContainerParams struct {
	ContainerName string `json:"container_name"`
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
}

func destroyContainer(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params destroyContainerParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, NewApiError(err, errBadInput, "could not decode data sent by client")
	}

	d := map[string]interface{}{
		"container_name":   ObjValue{params.ContainerName},
		"target_host":      ObjValue{params.TargetHost},
		"target_user_name": ObjValue{params.TargetUser},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, http.StatusBadRequest, NewApiError(err, errBadInput, "could not build actor's input")
	}

	return runSyncActor(req.Context(), "destroy-container", string(actorInput))
}
