package api

import (
	"encoding/json"
	"net/http"
)

type checkEntry struct {
	ID      string   `json:"check_id"`
	Status  string   `json:"status"`
	Summary string   `json:"summary"`
	Params  []string `json:"params"`
}

type checks struct {
	Checks []checkEntry `json:"checks"`
}

type checkHTMLOutputParams struct {
	CheckOutput []checks `json:"check_output"`
}

func checkHTMLOutputStart(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var params checkHTMLOutputParams

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		return nil, http.StatusBadRequest, newAPIError(err, errBadInput, "could not decode data sent by client")
	}

	j, err := json.Marshal(params)
	if err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "could not build actor's input")
	}
	id := actorRunnerRegistry.Create("check_html_output", string(j))

	s := actorRunnerRegistry.GetStatus(id, true)
	if err := checkSyncTaskStatus(s); err != nil {
		return nil, http.StatusInternalServerError, newAPIError(err, errInternal, "")
	}

	logExecutorError(req.Context(), s.Result)

	return parseExecutorResult(s.Result)
}
