package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPResponse func(http.ResponseWriter, *http.Request)

type ResultOk struct {
	Success bool
	Result  string
}

type ResultErr struct {
	Success bool
	ErrCode int
	Message string
}

func GenericHandler(f func(*json.Decoder) (string, error)) HTTPResponse {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)

		// TODO: as result return structures
		result, err := f(decoder)
		if err != nil {
			encoder.Encode("ERR.. :(")
		} else {
			encoder.Encode(result)
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
