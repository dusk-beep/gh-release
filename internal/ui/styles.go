package ui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true)

	metaStyle = lipgloss.NewStyle().
			Faint(true)

	separatorStyle = lipgloss.NewStyle().
			Faint(true)

	footerStyle = lipgloss.NewStyle().
			Faint(true)

	freshStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")) // green

	oldStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")) // orange

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Faint(true)
)
