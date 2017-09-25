package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPResponse func(http.ResponseWriter, *http.Request)

type ResultOk struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type ResultErr struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	Message string `json:"message"`
}

func GenericResponseHandler(f func(*http.Request) (interface{}, error)) HTTPResponse {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)

		result, err := f(request)
		if err != nil {
			// TODO: set appropriate err code
			// do we want to build our err structures with codes?
			encoder.Encode(ResultErr{false, 1, err.Error()})
		} else {
			encoder.Encode(ResultOk{true, result})
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
