package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type IssueStatus string

const (
	StatusTodo       IssueStatus = "Todo"
	StatusInProgress IssueStatus = "In Progress"
	StatusDone       IssueStatus = "Done"
)

type Issue struct {
	ID        int
	Title     string
	Status    IssueStatus
	CreatedAt time.Time
}

type Store struct {
	mu     sync.Mutex
	nextID int
	issues []Issue
}

func NewStore() *Store {
	return &Store{
		nextID: 1,
		issues: []Issue{
			{ID: 1, Title: "Create repo", Status: StatusDone, CreatedAt: time.Now().Add(-48 * time.Hour)},
			{ID: 2, Title: "Wire Gin + templates", Status: StatusInProgress, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{ID: 3, Title: "Add Postgres schema", Status: StatusTodo, CreatedAt: time.Now()},
		},
	}
}

func (s *Store) List() []Issue {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Issue, len(s.issues))
	copy(out, s.issues)
	return out
}

func (s *Store) Create(title string) Issue {
	s.mu.Lock()
	defer s.mu.Unlock()
	iss := Issue{
		ID:        s.nextID,
		Title:     title,
		Status:    StatusTodo,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.issues = append([]Issue{iss}, s.issues...)
	return iss
}

func (s *Store) SetStatus(id int, st IssueStatus) (Issue, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.issues {
		if s.issues[i].ID == id {
			s.issues[i].Status = st
			return s.issues[i], true
		}
	}
	return Issue{}, false
}

func main() {
	store := NewStore()

	r := gin.New()
	r.Use(gin.Recovery())

	// Templates + static
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "web/static")

	// Pages
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/issues")
	})

	r.GET("/issues", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":  "Taskusama — Issues",
			"Active": "issues",
			"Issues": store.List(),
		})
	})

	// HTMX: render only the <tbody> rows (used after create)
	r.GET("/issues/rows", func(c *gin.Context) {
		c.HTML(http.StatusOK, "issues.html", gin.H{
			"Issues": store.List(),
		})
	})

	// Create (HTMX)
	r.POST("/issues", func(c *gin.Context) {
		title := c.PostForm("title")
		if title == "" {
			c.String(http.StatusBadRequest, "title is required")
			return
		}
		issue := store.Create(title)

		// Return a single row fragment so HTMX can prepend it.
		c.HTML(http.StatusOK, "_issue_row.html", gin.H{
			"Issue": issue,
		})
	})

	// Update status (HTMX)
	r.POST("/issues/:id/status", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad id")
			return
		}
		status := IssueStatus(c.PostForm("status"))
		switch status {
		case StatusTodo, StatusInProgress, StatusDone:
		default:
			c.String(http.StatusBadRequest, "bad status")
			return
		}

		issue, ok := store.SetStatus(id, status)
		if !ok {
			c.String(http.StatusNotFound, "not found")
			return
		}

		// Return updated row HTML
		c.HTML(http.StatusOK, "_issue_row.html", gin.H{
			"Issue": issue,
		})
	})

	_ = r.Run(":8080")
}

