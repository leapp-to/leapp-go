package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPResponse func(http.ResponseWriter, *http.Request)

type ResultOk struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
}

type ResultErr struct {
	Success bool   `json:"success"`
	ErrCode int    `json:"err_code"`
	Message string `json:"message"`
}

func GenericHandler(f func(*json.Decoder) (map[string]interface{}, error)) HTTPResponse {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)

		result, err := f(decoder)
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
	router.HandleFunc("/migrate-machine", GenericHandler(MigrateMachine)).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
