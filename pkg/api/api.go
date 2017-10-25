package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

var actorRunnerRegistry *ActorRunner

func init() {
	// Initialize only once
	actorRunnerRegistry = NewActorRunner()
}

func parseExecutorResult(r *executor.Result) (interface{}, int, error) {
	if r.ExitCode != 0 {
		msg := fmt.Sprintf("actor execution failed with %d", r.ExitCode)
		return nil, http.StatusOK, NewApiError(nil, errActorExecution, msg)
	}

	if r.Stdout == "" {
		return nil, http.StatusOK, NewApiError(nil, errActorExecution, "actor didn't return any data")
	}

	var stdout interface{}
	if err := json.Unmarshal([]byte(r.Stdout), &stdout); err != nil {
		return nil, http.StatusOK, NewApiError(err, errActorExecution, "could not decode actor output")
	}
	return stdout, http.StatusOK, nil
}

// syncCmd wraps the result of the endpoint handler into a reponse that should be sent to the client.
func syncCmd(fn cmdFunc) respFunc {
	return func(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
		c, err := fn(req)
		if err != nil {
			return nil, http.StatusBadRequest, NewApiError(err, errBadInput, "error on endpoint handler execution")
		}

		r, err := c.Execute()
		// If requested, log actor's stderr
		v := req.Context().Value(CKey("Verbose")).(bool)
		if v {
			if err != nil {
				log.Printf("Actor execution failed: %s\n", err.Error())
			}
			log.Printf("Actor stderr: %s\n", r.Stderr)
		}

		return parseExecutorResult(r)
	}
}

func respHandler(fn respFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var r apiResult

		data, status, err := fn(rw, req)
		if err != nil {
			switch t := err.(type) {
			case apiError:
				r.Errors = append(r.Errors, t)
			default:
				http.Error(rw, "Internal error", http.StatusInternalServerError)
				return
			}
		} else {
			r.Data = data
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)

		err = json.NewEncoder(rw).Encode(r)
		if err != nil {
			log.Printf("could not encode response: %v\n", err)
		}
	}
}

// EndpointEntry represents an endpoint exposed by the daemon.
type EndpointEntry struct {
	Method      string
	Endpoint    string
	IsPrefix    bool
	NeedsStrip  bool
	HandlerFunc http.HandlerFunc
}

// GetEndpoints should return a slice of all endpoints that the daemon exposes.
func GetEndpoints() []EndpointEntry {
	return []EndpointEntry{
		{
			Method:      "POST",
			Endpoint:    "/migrate-machine",
			HandlerFunc: respHandler(migrateMachineStart),
		},
		{
			Method:      "GET",
			Endpoint:    "/migrate-machine/status/{id}",
			HandlerFunc: respHandler(migrateMachineResult),
		},
		{
			Method:      "POST",
			Endpoint:    "/port-inspect",
			HandlerFunc: respHandler(syncCmd(portInspectCmd)),
		},
		{
			Method:      "POST",
			Endpoint:    "/check-target",
			HandlerFunc: respHandler(syncCmd(checkTargetCmd)),
		},
		{
			Method:      "POST",
			Endpoint:    "/port-map",
			HandlerFunc: respHandler(syncCmd(portMapCmd)),
		},
		{
			Method:      "POST",
			Endpoint:    "/destroy-container",
			HandlerFunc: respHandler(syncCmd(destroyContainerCmd)),
		},
		{
			Method:      "GET",
			Endpoint:    "/doc",
			IsPrefix:    true,
			NeedsStrip:  true,
			HandlerFunc: http.FileServer(http.Dir("/usr/share/leapp/apidoc/")).ServeHTTP,
		},
	}
}
