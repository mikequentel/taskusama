package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mikequentel/taskusama/internal/model"
	"github.com/mikequentel/taskusama/internal/store"
)

type Server struct {
	issues store.IssueStore
}

func New(issues store.IssueStore) *Server {
	return &Server{issues: issues}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	// pages
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/issues")
	})

	r.GET("/issues", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":  "Taskusama — Issues",
			"Active": "issues",
			"Issues": s.issues.List(),
		})
	})

	// HTMX: render only the <tbody> rows (used after create)
	r.GET("/issues/rows", func(c *gin.Context) {
		c.HTML(http.StatusOK, "issues.html", gin.H{
			"Issues": s.issues.List(),
		})
	})

	// creates (HTMX)
	r.POST("/issues", func(c *gin.Context) {
		title := c.PostForm("title")
		if title == "" {
			c.String(http.StatusBadRequest, "title is required")
			return
		}
		issue := s.issues.Create(title)

		// IMPORTANT: your _issue_row.html expects dot = Issue,
		// so pass the Issue directly, not gin.H{"Issue": issue}.
		c.HTML(http.StatusOK, "_issue_row.html", issue)
	})

	// updates status (HTMX)
	r.POST("/issues/:id/status", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "bad id")
			return
		}

		status := model.IssueStatus(c.PostForm("status"))
		switch status {
		case model.StatusTodo, model.StatusInProgress, model.StatusDone:
		default:
			c.String(http.StatusBadRequest, "bad status")
			return
		}

		issue, ok := s.issues.SetStatus(id, status)
		if !ok {
			c.String(http.StatusNotFound, "not found")
			return
		}

		c.HTML(http.StatusOK, "_issue_row.html", issue)
	})
}
