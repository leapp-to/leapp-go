package api

import (
	"encoding/json"
	"net/http"
)

type checkTargetParams struct {
	TargetHost string `json:"target_host"`
	Status     bool   `json:"check_target_service_status"`
	TargetUser string `json:"target_user_name"`
}

func buildCheckTargetInput(p *checkTargetParams) (string, error) {
	data := map[string]interface{}{
		"target_host":                 ObjValue{p.TargetHost},
		"check_target_service_status": ObjValue{p.Status},
		"target_user_name":            ObjValue{p.TargetUser},
	}

	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func checkTarget(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	if err := checkSyncTaskStatus(nil); err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "")
	}
	var params checkTargetParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	actorInput, err := buildCheckTargetInput(&params)
	if err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "could not build actor's input")
	}

	id := actorRunnerRegistry.Create("remote-target-check-group", actorInput)

	s := actorRunnerRegistry.GetStatus(id, true)
	if err := checkSyncTaskStatus(s); err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)

}
