package api

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse func(http.ResponseWriter, *http.Request)

type Result struct {
	Output     interface{} `json:"output,omitempty"`
	ErrCode    int         `json:"err_code"`
	ErrMessage string      `json:"err_message,omitempty"`
}

func GenericResponseHandler(fn func(*http.Request) (interface{}, error)) HTTPResponse {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)

		result, err := fn(request)
		if err != nil {
			// TODO: set appropriate err code
			// do we want to build our err structures with codes?
			encoder.Encode(Result{ErrCode: 1, ErrMessage: err.Error()})
		} else {
			encoder.Encode(Result{ErrCode: 0, Output: result})
		}
	}
}

type EndpointEntry struct {
	Method      string
	Endpoint    string
	HandlerFunc HTTPResponse
}

func GetEndpoints() []EndpointEntry {
	return []EndpointEntry{
		{
			Method:      "POST",
			Endpoint:    "/migrate-machine",
			HandlerFunc: GenericResponseHandler(MigrateMachine),
		},
	}
}
