package ui

import (
	"fmt"
	"gh-release/internal/model"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchCmd(),
	)
}

func fmtReleaseTitle(release model.Release) string {
	if release.Name != "" && release.Name != release.Tag {
		return release.Tag + " — " + release.Name
	}

	return release.Tag
}

func (m Model) openSelectedReleaseCmd() tea.Cmd {
	return func() tea.Msg {
		row := m.table.SelectedRow()

		for _, result := range m.results {
			if result.Err != nil {
				continue
			}

			if result.Repository.FullName() == row[0] {
				rendered, err := renderReleaseNotes(
					result.Release,
					max(20, m.width-4),
				)

				return releaseNotesMsg{
					title: fmtReleaseTitle(result.Release),
					body:  rendered,
					err:   err,
				}
			}
		}

		return releaseNotesMsg{
			err: fmt.Errorf("release not found"),
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		m.resizeTable()

	case tea.KeyPressMsg:
		if m.showRelease {
			switch msg.String() {
			case "esc", "q":
				m.showRelease = false
				m.releaseScroll = 0

			case "up", "k":
				m.releaseScroll = max(0, m.releaseScroll-1)

			case "down", "j":
				m.releaseScroll++

			case "pgup":
				m.releaseScroll = max(0, m.releaseScroll-(m.height/2))

			case "pgdown", "space":
				m.releaseScroll += m.height / 2

			case "home", "g":
				m.releaseScroll = 0

			case "end", "G":
				m.releaseScroll = max(0, m.releaseScroll+100000)
			}

			return m, nil
		}

		switch msg.String() {
		case "enter":
			return m, m.openSelectedReleaseCmd()

		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			if !m.loading {
				m.loading = true
				cmds = append(cmds, m.fetchCmd())
			}
		}

	case fetchDoneMsg:
		m.loading = false
		m.results = msg.results
		m.rebuildTable()

	case releaseNotesMsg:
		if msg.err != nil {
			return m, nil
		}

		m.showRelease = true
		m.releaseTitle = msg.title
		m.releaseMarkdown = msg.body
		m.releaseScroll = 0

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) resize(width, height int) {
	m.width = width
	m.height = height

	m.table.SetWidth(max(20, width-4))
	m.table.SetHeight(max(5, height-8))
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
