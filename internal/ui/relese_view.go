package ui

import (
	"fmt"
	"strings"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"

	"gh-release/internal/model"
)

func renderReleaseNotes(release model.Release, width int) (string, error) {
	content := release.Body

	if strings.TrimSpace(content) == "" {
		content = "_No release notes._"
	}

	renderWidth := max(20, width-4)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(renderWidth),
		glamour.WithTableWrap(false),
	)
	if err != nil {
		return "", fmt.Errorf("create markdown renderer: %w", err)
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return "", fmt.Errorf("render release notes: %w", err)
	}

	title := release.Tag
	if release.Name != "" && release.Name != release.Tag {
		title += " — " + release.Name
	}

	header := titleStyle.Render(title)

	return header + "\n\n" + strings.TrimRight(rendered, "\n"), nil
}

func releaseNotesView(title, content string, width, height, scroll int) string {
	lines := strings.Split(content, "\n")

	contentWidth := max(20, width-4)
	contentHeight := max(1, height-5)

	if scroll > len(lines)-contentHeight {
		scroll = max(0, len(lines)-contentHeight)
	}

	end := min(len(lines), scroll+contentHeight)

	var b strings.Builder

	b.WriteString("╭")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("╮")
	b.WriteString("\n")

	header := titleStyle.Render(title) +
		"  " +
		footerStyle.Render("esc back  •  ↑↓ scroll")

	writeReleaseLine(&b, header, contentWidth)

	b.WriteString("│ ")
	b.WriteString(strings.Repeat("─", contentWidth))
	b.WriteString(" │\n")

	for _, line := range lines[scroll:end] {
		writeReleaseLine(&b, line, contentWidth)
	}

	for i := end - scroll; i < contentHeight; i++ {
		writeReleaseLine(&b, "", contentWidth)
	}

	b.WriteString("╰")
	b.WriteString(strings.Repeat("─", width-2))
	b.WriteString("╯")

	return b.String()
}

func writeReleaseLine(b *strings.Builder, line string, width int) {
	lineWidth := lipgloss.Width(line)

	b.WriteString("│ ")
	b.WriteString(line)

	if lineWidth < width {
		b.WriteString(strings.Repeat(" ", width-lineWidth))
	}

	b.WriteString(" │\n")
}
