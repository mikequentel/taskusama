package model

import "time"

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
