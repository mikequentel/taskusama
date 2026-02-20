package memory

import (
	"sync"
	"time"

	"github.com/mikequentel/taskusama/internal/model"
)

type Store struct {
	mu     sync.Mutex
	nextID int
	issues []model.Issue
}

func New() *Store {
	now := time.Now()
	return &Store{
		nextID: 4,
		issues: []model.Issue{
			{ID: 1, Title: "Create repo", Status: model.StatusDone, CreatedAt: now.Add(-48 * time.Hour)},
			{ID: 2, Title: "Wire Gin + templates", Status: model.StatusInProgress, CreatedAt: now.Add(-2 * time.Hour)},
			{ID: 3, Title: "Add Postgres schema", Status: model.StatusTodo, CreatedAt: now},
		},
	}
}

func (s *Store) List() []model.Issue {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.Issue, len(s.issues))
	copy(out, s.issues)
	return out
}

func (s *Store) Create(title string) model.Issue {
	s.mu.Lock()
	defer s.mu.Unlock()

	iss := model.Issue{
		ID:        s.nextID,
		Title:     title,
		Status:    model.StatusTodo,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.issues = append([]model.Issue{iss}, s.issues...)
	return iss
}

func (s *Store) SetStatus(id int, st model.IssueStatus) (model.Issue, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.issues {
		if s.issues[i].ID == id {
			s.issues[i].Status = st
			return s.issues[i], true
		}
	}
	return model.Issue{}, false
}
