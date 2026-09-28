package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CasimirG12/jobrunner/internal/job"
	"github.com/CasimirG12/jobrunner/internal/store"
)

const (
	maxBodyBytes = 1 << 20
)

func writeJSON(w http.ResponseWriter, v any, code int) {
	body, err := json.Marshal(v)

	if err != nil {
		log.Printf("marshal response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	errorRes := errorResponse{
		Error: message,
	}
	writeJSON(w, errorRes, code)
}

type CreateJobReq struct {
	Kind  job.Kind `json:"kind"`
	Input string   `json:"input"`
}

type Server struct {
	store *store.MemoryStore
}

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)

	var body CreateJobReq
	err := decoder.Decode(&body)
	if err != nil {
		var maxByteError *http.MaxBytesError
		if errors.As(err, &maxByteError) {
			writeError(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}
		writeError(w, "bad request", http.StatusBadRequest)
		return
	}

	newJob, err := job.New(body.Input, body.Kind)
	if err != nil {
		if errors.Is(err, job.ErrInvalidKind) {
			writeError(w, "job kind does not exist", http.StatusBadRequest)
			return
		}
		log.Printf("create job: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.store.Create(newJob)
	if err != nil {
		log.Printf("create job: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, newJob, http.StatusCreated)
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.store.GetByID(id)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, "job not found", http.StatusNotFound)
			return
		}
		log.Printf("get job: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, result, http.StatusOK)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	list, err := s.store.List()
	if err != nil {
		log.Printf("list jobs: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, list, http.StatusOK)
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", s.createJob)
	mux.HandleFunc("GET /jobs", s.listJobs)
	mux.HandleFunc("GET /jobs/{id}", s.getJob)
	return mux
}

func NewServer() *Server {
	return &Server{
		store: store.NewMemoryStore(),
	}
}
