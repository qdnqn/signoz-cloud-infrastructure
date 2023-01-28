package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-backend/internal"
	"go-backend/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net/http"
)

type handler struct {
	Service internal.Service
	config  config.Config
}

func Handler(s internal.Service, c config.Config) http.Handler {
	h := &handler{s, c}
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("DemoService"))

	r.HandleFunc("/putEntry", h.CreateOrUpdateEntry).Methods("PUT")
	r.HandleFunc("/getEntries", h.GetEntries).Methods("GET")
	r.HandleFunc("/getKey/{key}", h.GetKey).Methods("GET")
	r.HandleFunc("/getEntries/{id}", h.GetEntry).Methods("GET")
	r.HandleFunc("/deleteEntry/{id}", h.DeleteEntry).Methods("DELETE")

	return r
}

func (h *handler) CreateOrUpdateEntry(w http.ResponseWriter, r *http.Request) {
	var entry internal.Entry
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	out, err := h.Service.CreateOrUpdateEntry(ctx, entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.config.Logger.Info("Emiting event for entry id: " + out.Id)

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *handler) GetEntries(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()
	out, err := h.Service.GetEntries(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *handler) GetKey(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	ctx := r.Context()
	out, err := h.Service.GetKey(ctx, vars["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	ctx := r.Context()
	_, span := h.config.Tracer.Start(ctx, "getEntry", oteltrace.WithAttributes(attribute.String("id", vars["id"])))
	defer span.End()

	h.config.Logger.Info("Trying to get entry: " + vars["id"])
	out, err := h.Service.GetEntry(ctx, vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	ctx := r.Context()
	err := h.Service.DeleteEntry(ctx, vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
