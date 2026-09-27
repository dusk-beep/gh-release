package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"gh-release/internal/service"
)

func (m *Model) startFetch() tea.Cmd {
	if m.fetchCancel != nil {
		m.fetchCancel()
	}

	m.fetchCtx, m.fetchCancel = context.WithCancel(context.Background())

	m.fetchResults = m.service.Fetch(
		m.fetchCtx,
		m.repos,
	)

	return waitForResult(m.fetchResults)
}

func waitForResult(ch <-chan service.Result) tea.Cmd {
	return func() tea.Msg {
		result, ok := <-ch

		if !ok {
			return fetchCompleteMsg{}
		}

		return repositoryResultMsg{
			index:  result.Index,
			result: result.Result,
		}
	}
}
