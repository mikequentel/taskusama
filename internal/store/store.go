package store

import "github.com/mikequentel/taskusama/internal/model"

type IssueStore interface {
	List() []model.Issue
	Create(title string) model.Issue
	SetStatus(id int, st model.IssueStatus) (model.Issue, bool)
}
