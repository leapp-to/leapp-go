package actors

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// ActorAPIService is the implementation auf the audit service
type ActorAPIService struct {
	stop   chan struct{}
	server http.Server
	router *mux.Router
	dao    ActorDAO
}

func getenv(name, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

// NewActorAPIService initializes the instance of the audit service
func NewActorAPIService(block bool) (service *ActorAPIService, err error) {
	dbx, err := sqlx.Open("sqlite3", getenv("LEAPP_STORE_PATH", "/var/lib/leapp/leapp-store.db"))
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	sub := router.PathPrefix("/actors/v1/").Subrouter()

	service = &ActorAPIService{
		stop:   make(chan struct{}),
		router: router,
		server: http.Server{Handler: sub},
		dao:    ActorDAO{dbx},
	}

	sub.HandleFunc("/audit", func(w http.ResponseWriter, r *http.Request) {
		var audit Audit
		if err := json.NewDecoder(r.Body).Decode(&audit); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}).Methods("POST")

	sub.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		var logMsg LogMessage
		if err := json.NewDecoder(r.Body).Decode(&logMsg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		service.dao.AddAudit(convertLogMessage(logMsg))
	}).Methods("POST")

	sub.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		msg := new(Message)
		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		service.dao.AddAudit(convertMessage(msg))
	}).Methods("POST")

	sub.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		var query struct {
			Context  string
			Messages []string
		}
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msgs := []Message{}
		if len(query.Messages) > 0 {
			msgs, err = service.dao.GetMessages(query.Context, query.Messages)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				fmt.Printf("QUERY FAILURE: %s\n", err.Error())
				return
			}
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Messages []Message `json:"messages"`
		}{msgs})
	}).Methods("POST")

	socketPath := getenv("LEAPP_ACTOR_API", "/var/run/leapp-actor-api.sock")
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	go service.run()
	if block {
		service.server.Serve(listener)
	} else {
		go service.server.Serve(listener)
	}

	return service, nil
}

// Close stops the audit service
func (service *ActorAPIService) Close() error {
	if service != nil {
		err := service.server.Close()
		service.stop <- struct{}{}
		return err
	}
	return nil
}

func (service *ActorAPIService) run() {
	for {
		select {
		case <-service.stop:
			break
		}
	}
}
