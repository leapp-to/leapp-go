package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leapp-to/leapp-go/pkg/msg"
)

func MigrateHandler(w http.ResponseWriter, r *http.Request) {
	var migrateParams msg.MigrateParams
	encoder := json.NewEncoder(w)

	err := json.NewDecoder(r.Body).Decode(&migrateParams)
	if err != nil {
		// replace that with logging error
		fmt.Println(err)

		// build errors
		encoder.Encode(msg.Result{false, "So sad.. :( error occured"})
	} else {
		// validate params

		// send params to executor

		// send it back to ui
		json.NewEncoder(w).Encode(&migrateParams)
	}

}

func RunHTTPServer() {

	router := mux.NewRouter()
	router.HandleFunc("/migrate-machine", MigrateHandler).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
