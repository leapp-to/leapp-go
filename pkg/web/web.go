package web

import (
	"context"
	"log"
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

	Verbose bool
}

// Handler contains everything needed to start the HTTP service.
type Handler struct {
	options *Options
	mux     *mux.Router
	errorCh chan error
}

// ServeHTTP implements the Handler interface.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if h.options.Verbose {
		log.Printf("\"%s %s\" - %s", req.Method, req.RequestURI, req.RemoteAddr)
	}

	ctx := context.WithValue(context.Background(), "Verbose", h.options.Verbose)
	h.mux.ServeHTTP(rw, req.WithContext(ctx))
}

// Run serves the HTTP endpoints.
func (h *Handler) Run() {
	srv := &http.Server{
		Handler:      h,
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
