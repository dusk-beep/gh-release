package ui

import (
	"sort"
	"time"

	"charm.land/bubbles/v2/table"

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

	styles.Selected = styles.Selected.
		Bold(true)

	t.SetStyles(styles)

	return t
}

func (m *Model) rebuildTable() {
	results := append([]model.RepositoryResult(nil), m.results...)

	sort.SliceStable(results, func(i, j int) bool {
		a := results[i].Release.PublishedAt
		b := results[j].Release.PublishedAt

		// Releases with no date go last.
		if a.IsZero() {
			return false
		}

		if b.IsZero() {
			return true
		}

		return a.After(b)
	})

	rows := make([]table.Row, 0, len(results))

	for _, result := range results {
		rows = append(rows, resultRow(result))
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
		release.Repository.FullName(),
		releaseName(release),
		formatPublished(release.PublishedAt),
		releaseStatus(release),
	}
}

func releaseName(release model.Release) string {
	if release.Name == "" || release.Name == release.Tag {
		return release.Tag
	}

	return release.Tag + " · " + release.Name
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
		return "draft"
	case release.Prerelease:
		return "pre-release"
	default:
		return "release"
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

	repositoryWidth := width * 35 / 100
	releaseWidth := width * 30 / 100
	publishedWidth := width * 20 / 100
	statusWidth := width * 15 / 100

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
