package model

type IssueID int

type Issue struct {
	Description string
}

func NewIssue(description string) *Issue {
	return &Issue{
		Description: description,
	}
}
