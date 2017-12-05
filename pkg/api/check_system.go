package api

import (
	"encoding/json"
	"net/http"
)

type checkSystemParams struct {
	Checks string `json:"checks"`
}

func checkSystemStart(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params checkSystemParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	id := actorRunnerRegistry.Create("check_"+params.Checks, "{}")

	s := actorRunnerRegistry.GetStatus(id, true)
	if err := checkSyncTaskStatus(s); err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}
