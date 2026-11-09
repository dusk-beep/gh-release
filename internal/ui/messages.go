package ui

import "gh-release/internal/model"

type fetchDoneMsg struct {
	results []model.RepositoryResult
}

type releaseNotesMsg struct {
	err   error
	title string
	body  string
}

type repositoryResultMsg struct {
	result model.RepositoryResult
	index  int
}

type fetchCompleteMsg struct{}
