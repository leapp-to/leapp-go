package web

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leapp-to/leapp-go/pkg/api"
)

// Options contains parameters for the web handler.
type Options struct {
	ListenAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// Handler contains everything needed to start the HTTP service.
type Handler struct {
	options *Options
	mux     *mux.Router
	errorCh chan error
}

// Run serves the HTTP endpoints.
func (h *Handler) Run() {
	srv := &http.Server{
		Handler:      h.mux,
		Addr:         h.options.ListenAddress,
		ReadTimeout:  h.options.ReadTimeout * time.Second,
		WriteTimeout: h.options.WriteTimeout * time.Second,
	}

	if listener, err := net.Listen("tcp", srv.Addr); err == nil {
		h.errorCh <- srv.Serve(listener)
	} else {
		h.errorCh <- err
	}
}

// ErrorCh returns a channel where the web handler errors go to.
func (h *Handler) ErrorCh() <-chan error {
	return h.errorCh
}

// New initializes a new Handler.
func New(o *Options) *Handler {
	h := &Handler{
		mux:     mux.NewRouter(),
		errorCh: make(chan error),
		options: o,
	}

	apiV1 := h.mux.PathPrefix("/v1").Subrouter()

	for _, e := range api.GetEndpoints() {
		apiV1.HandleFunc(e.Endpoint, e.HandlerFunc).Methods(e.Method)
	}

	return h
}
