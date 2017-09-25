package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

func RunHTTPServer() {
	router := mux.NewRouter()
	apiV1 := router.PathPrefix("/v1.0").Subrouter()
	apiV1.HandleFunc("/migrate-machine", GenericResponseHandler(MigrateMachine)).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
