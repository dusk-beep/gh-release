package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
)

func (m Model) fetchCmd() tea.Cmd {
	return func() tea.Msg {
		results := m.service.Fetch(
			context.Background(),
			m.repos,
		)

		return fetchDoneMsg{
			results: results,
		}
	}
}
