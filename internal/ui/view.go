package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if m.showRelease {
		return m.releaseView()
	}

	var b strings.Builder

	width := max(60, m.width)

	// Top border.
	b.WriteString("╭")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("╮")
	b.WriteString("\n")

	// Header.
	b.WriteString("│ ")
	m.renderHeader(&b)
	b.WriteString(strings.Repeat(" ", max(0, width-3-lipgloss.Width(b.String()))))
	b.WriteString("│")
	b.WriteString("\n")

	// Separator.
	b.WriteString("│ ")
	b.WriteString(strings.Repeat("─", width-4))
	b.WriteString(" │")
	b.WriteString("\n")

	// Table.
	tableWidth := max(20, width-4)
	m.table.SetWidth(tableWidth)

	tableView := m.table.View()

	for _, line := range strings.Split(tableView, "\n") {
		b.WriteString("│ ")
		b.WriteString(line)
		padding := tableWidth - lipgloss.Width(line)
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString(" │\n")
	}

	// Footer.
	b.WriteString("│ ")
	m.renderFooter(&b)
	b.WriteString(" │\n")

	// Bottom border.
	b.WriteString("╰")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("╯")

	view := tea.NewView(b.String())
	view.AltScreen = true

	return view
}

func (m Model) renderHeader(b *strings.Builder) {
	header := titleStyle.Render("gh-release") +
		"  " +
		metaStyle.Render(fmt.Sprintf("%d repositories", len(m.repos)))

	if m.loading {
		header += "  " +
			m.spinner.View() +
			" " +
			metaStyle.Render("fetching")
	}

	b.WriteString(header)

	padding := max(0, m.width-4-lipgloss.Width(header))
	b.WriteString(strings.Repeat(" ", padding))
}

func (m Model) renderFooter(b *strings.Builder) {
	footer := footerStyle.Render("r refresh  •  q quit")

	b.WriteString(footer)

	padding := max(0, m.width-4-lipgloss.Width(footer))
	b.WriteString(strings.Repeat(" ", padding))
}

func formatRelative(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Second:
		return "just now"
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func (m Model) releaseView() tea.View {
	width := max(60, m.width)
	height := max(10, m.height)

	view := releaseNotesView(
		m.releaseTitle,
		m.releaseMarkdown,
		width,
		height,
		m.releaseScroll,
	)

	v := tea.NewView(view)
	v.AltScreen = true

	return v
}
