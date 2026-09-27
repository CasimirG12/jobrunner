package store

import (
	"errors"
	"sync"

	"github.com/CasimirG12/jobrunner/internal/job"
)

var ErrNotFound = errors.New("store: job not found")
var ErrAlreadyExists = errors.New("store: job id already exists")

type MemoryStore struct {
	jobs map[string]job.Job
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]job.Job),
	}
}

func (s *MemoryStore) GetByID(id string) (job.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.jobs[id]
	if !ok {
		return job.Job{}, ErrNotFound
	}
	return result, nil
}

func (s *MemoryStore) List() ([]job.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]job.Job, 0, len(s.jobs))

	for _, j := range s.jobs {
		result = append(result, j)
	}
	return result, nil
}

func (s *MemoryStore) Create(j job.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[j.ID]; ok {
		return ErrAlreadyExists
	}
	s.jobs[j.ID] = j
	return nil
}
