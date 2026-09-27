package ui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	repositoryOwnerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7"))

	repositorySeparatorStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("8"))

	repositoryNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("14")).
				Bold(true)

	releaseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)

	freshStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	oldStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	statusReleaseStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("14"))

	statusPrereleaseStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11"))

	statusDraftStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	fetchingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))
)
