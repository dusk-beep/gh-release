package ui

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"

	"gh-release/internal/model"
	"gh-release/internal/service"
)

type Model struct {
	service *service.Service
	repos   []model.Repository

	results      []model.RepositoryResult
	fetchResults <-chan service.Result

	table   table.Model
	spinner spinner.Model

	showRelease     bool
	releaseMarkdown string
	releaseTitle    string
	releaseScroll   int

	loading bool

	width  int
	height int
}

func New(svc *service.Service, repos []model.Repository) Model {
	results := make([]model.RepositoryResult, len(repos))

	for i, repo := range repos {
		results[i] = model.RepositoryResult{
			Repository: repo,
		}
	}

	return Model{
		service: svc,
		repos:   repos,
		results: results,
		table:   newTable(),
		spinner: spinner.New(),
		loading: true,
	}
}
