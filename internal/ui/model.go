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

	results []model.RepositoryResult

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
	return Model{
		service: svc,
		repos:   repos,
		table:   newTable(),
		spinner: spinner.New(),
		loading: true,
	}
}
