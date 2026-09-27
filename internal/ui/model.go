package ui

import (
	"context"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"

	"gh-release/internal/model"
	"gh-release/internal/service"
)

type Model struct {
	service *service.Service
	repos   []model.Repository

	results []model.RepositoryResult
	table   table.Model
	spinner spinner.Model

	loading bool

	fetchCtx     context.Context
	fetchCancel  context.CancelFunc
	fetchResults <-chan service.Result

	width  int
	height int

	displayOrder []int

	showRelease     bool
	releaseMarkdown string
	releaseTitle    string
	releaseScroll   int
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
