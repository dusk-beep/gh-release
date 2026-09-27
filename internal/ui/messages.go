package ui

import "gh-release/internal/model"

type fetchDoneMsg struct {
	results []model.RepositoryResult
}

type releaseNotesMsg struct {
	title string
	body  string
	err   error
}

type repositoryResultMsg struct {
	index  int
	result model.RepositoryResult
}

type fetchCompleteMsg struct{}
