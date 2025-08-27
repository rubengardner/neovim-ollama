package reviews

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rubengardner/neovim-ollama/cmd/neovim-ollama/ui"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func RenderReviewChanges(m *model.Model) string {
	styles := ui.NewStyles()
	if len(m.ProposedChanges) == 0 {
		return "No changes to review"
	}

	// Ensure cursor is within bounds
	if m.ReviewCursor < 0 {
		m.ReviewCursor = 0
	}
	if m.ReviewCursor >= len(m.ProposedChanges) {
		m.ReviewCursor = len(m.ProposedChanges) - 1
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("Review Changes (%d total)\n\n", len(m.ProposedChanges)))

	for i, change := range m.ProposedChanges {
		prefix := "  "
		if i == m.ReviewCursor {
			prefix = "> "
		}

		// Status indicator
		status := "⏳"                                                            // pending
		statusColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")) // orange
		if change.Accepted != nil {
			if *change.Accepted {
				status = "✅"                                                            // accepted
				statusColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // green
			} else {
				status = "❌"                                                            // rejected
				statusColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")) // red
			}
		}

		// File and line info
		fileName := filepath.Base(change.FilePath)
		lineInfo := ""
		if change.LineStart > 0 {
			if change.LineEnd > change.LineStart {
				lineInfo = fmt.Sprintf(" (lines %d-%d)", change.LineStart, change.LineEnd)
			} else {
				lineInfo = fmt.Sprintf(" (line %d)", change.LineStart)
			}
		}

		line := fmt.Sprintf("%s%s %s%s", prefix, statusColor.Render(status), fileName, lineInfo)

		if i == m.ReviewCursor {
			line = styles.Selected.Render(line)
		}

		content.WriteString(line + "\n")

		// Show description if this is the current item
		if i == m.ReviewCursor && change.Description != "" {
			desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true).Render("  " + change.Description)
			content.WriteString(desc + "\n")
		}

		// Show diff preview if this is the current item
		if i == m.ReviewCursor {
			content.WriteString("\n  --- Original ---\n")
			originalLines := strings.Split(change.OriginalCode, "\n")
			for _, originalLine := range originalLines {
				if strings.TrimSpace(originalLine) != "" {
					content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("  - "+originalLine) + "\n")
				}
			}

			content.WriteString("\n  +++ Proposed +++\n")
			proposedLines := strings.Split(change.ProposedCode, "\n")
			for _, proposedLine := range proposedLines {
				if strings.TrimSpace(proposedLine) != "" {
					content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#51CF66")).Render("  + "+proposedLine) + "\n")
				}
			}
			content.WriteString("\n")
		}
	}

	// Summary
	accepted := 0
	rejected := 0
	pending := 0
	for _, change := range m.ProposedChanges {
		if change.Accepted == nil {
			pending++
		} else if *change.Accepted {
			accepted++
		} else {
			rejected++
		}
	}

	content.WriteString(fmt.Sprintf("\nStatus: ✅ %d accepted, ❌ %d rejected, ⏳ %d pending", accepted, rejected, pending))

	return content.String()
}
