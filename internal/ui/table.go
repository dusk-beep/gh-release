package ui

import (
	"sort"
	"time"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"

	"gh-release/internal/model"
)

func newTable() table.Model {
	columns := []table.Column{
		{
			Title: "Repository",
			Width: 36,
		},
		{
			Title: "Latest release",
			Width: 30,
		},
		{
			Title: "Published",
			Width: 20,
		},
		{
			Title: "Status",
			Width: 14,
		},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	styles := table.DefaultStyles()

	styles.Header = styles.Header.
		Bold(true)

	styles.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("237"))

	t.SetStyles(styles)

	return t
}

func (m *Model) rebuildTable() {
	indices := make([]int, len(m.results))

	for i := range indices {
		indices[i] = i
	}

	sort.SliceStable(indices, func(i, j int) bool {
		a := m.results[indices[i]].Release.PublishedAt
		b := m.results[indices[j]].Release.PublishedAt

		if a.IsZero() {
			return false
		}

		if b.IsZero() {
			return true
		}

		return a.After(b)
	})

	m.displayOrder = indices

	rows := make([]table.Row, 0, len(indices))

	for _, index := range indices {
		rows = append(rows, resultRow(m.results[index]))
	}

	m.table.SetRows(rows)
}

func resultRow(result model.RepositoryResult) table.Row {
	if result.Err != nil {
		return table.Row{
			result.Repository.FullName(),
			"",
			"",
			errorStyle.Render(shortError(result.Err)),
		}
	}

	release := result.Release

	return table.Row{
		formatRepository(release.Repository),
		releaseName(release),
		formatPublished(release.PublishedAt),
		releaseStatus(release),
	}
}

func formatRepository(repo model.Repository) string {
	return repositoryOwnerStyle.Render(repo.Owner) +
		repositorySeparatorStyle.Render("/") +
		repositoryNameStyle.Render(repo.Name)
}

func releaseName(release model.Release) string {
	return releaseStyle.Render(release.Tag)
}

func formatPublished(t time.Time) string {
	if t.IsZero() {
		return "-"
	}

	age := time.Since(t)

	text := formatRelative(t)

	if age < 24*time.Hour {
		return freshStyle.Render(text)
	}

	return oldStyle.Render(text)
}

func releaseStatus(release model.Release) string {
	switch {
	case release.Draft:
		return statusDraftStyle.Render("draft")

	case release.Prerelease:
		return statusPrereleaseStyle.Render("pre-release")

	default:
		return statusReleaseStyle.Render("release")
	}
}

func shortError(err error) string {
	const maxLength = 24

	s := err.Error()

	if len(s) <= maxLength {
		return s
	}

	return s[:maxLength-3] + "..."
}

func (m *Model) resizeTable() {
	width := max(60, m.width-6)

	repositoryWidth := width * 40 / 100
	releaseWidth := width * 30 / 100
	publishedWidth := width * 18 / 100
	statusWidth := width - repositoryWidth - releaseWidth - publishedWidth

	m.table.SetColumns([]table.Column{
		{
			Title: "Repository",
			Width: repositoryWidth,
		},
		{
			Title: "Latest release",
			Width: releaseWidth,
		},
		{
			Title: "Published",
			Width: publishedWidth,
		},
		{
			Title: "Status",
			Width: statusWidth,
		},
	})
}
