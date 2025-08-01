package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	ChatMode Mode = iota
	FileSelectMode
)

type FileItem struct {
	Name     string
	Path     string
	IsDir    bool
	Selected bool
}

type Model struct {
	input         textinput.Model
	viewport      viewport.Model
	width         int
	height        int
	isWaiting     bool
	err           error
	spinner       spinner.Model
	history       []ChatMessage
	mode          Mode
	files         []FileItem
	filesCursor   int
	currentDir    string
	selectedFiles []string
	filesViewport viewport.Model
}

type (
	responseMsg    string
	errorMsg       error
	filesLoadedMsg []FileItem
)

var (
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63"))
	inputStyle    = borderStyle.Padding(0, 1)
	outputStyle   = borderStyle.Padding(0, 1).MarginBottom(1)
	promptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Bold(true)
	responseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F"))
	divider       = lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(strings.Repeat("─", 40))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#3C3836")).Foreground(lipgloss.Color("#EBDBB2"))
	fileStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#83A598"))
	folderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FABD2F")).Bold(true)
	checkedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26")).Bold(true)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")).Italic(true)
)

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter prompt (Ctrl+F for files, Ctrl+C to exit)"
	ti.Focus()
	ti.CharLimit = 500

	vp := viewport.New(0, 0)

	filesVp := viewport.New(0, 0)

	sp := spinner.New()
	sp.Spinner = spinner.Line
	sp.Style = lipgloss.NewStyle()

	currentDir, _ := os.Getwd()

	return Model{
		input:         ti,
		viewport:      vp,
		filesViewport: filesVp,
		spinner:       sp,
		mode:          ChatMode,
		currentDir:    currentDir,
		selectedFiles: []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func loadFiles(dir string) tea.Cmd {
	return func() tea.Msg {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return errorMsg(err)
		}

		var fileItems []FileItem

		// Add parent directory if not at root
		if dir != "/" && dir != "." {
			fileItems = append(fileItems, FileItem{
				Name:  "..",
				Path:  filepath.Dir(dir),
				IsDir: true,
			})
		}

		for _, file := range files {
			// Skip hidden files and directories
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			fileItems = append(fileItems, FileItem{
				Name:  file.Name(),
				Path:  filepath.Join(dir, file.Name()),
				IsDir: file.IsDir(),
			})
		}

		return filesLoadedMsg(fileItems)
	}
}

func readFileContent(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = m.width - 4
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 7 // Leave room for help text
		m.viewport.YPosition = 1
		m.filesViewport.Width = m.width - 4
		m.filesViewport.Height = m.height - 7
		m.filesViewport.YPosition = 1

		if m.mode == ChatMode {
			m.viewport.SetContent(renderHistory(m.history))
		} else {
			m.filesViewport.SetContent(m.renderFiles())
		}
		return m, nil

	case tea.KeyMsg:
		if m.mode == ChatMode {
			return m.handleChatKeys(msg, &cmds)
		} else {
			return m.handleFileKeys(msg, &cmds)
		}

	case responseMsg:
		m.isWaiting = false
		m.input.SetValue("")
		m.input.CursorEnd()
		m.input.Focus()
		m.history[len(m.history)-1].Content = string(msg)
		m.viewport.SetContent(renderHistory(m.history))
		return m, nil

	case errorMsg:
		m.isWaiting = false
		m.input.SetValue("")
		m.input.Focus()
		if len(m.history) > 0 {
			m.history[len(m.history)-1].Content = "Error: " + msg.Error()
			m.viewport.SetContent(renderHistory(m.history))
		}
		return m, nil

	case filesLoadedMsg:
		m.files = []FileItem(msg)
		m.filesCursor = 0
		m.filesViewport.SetContent(m.renderFiles())
		return m, nil
	}

	if m.isWaiting {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinCmd)
		m.input.SetValue(m.spinner.View() + " Thinking...")
		m.input.SetCursor(0)
	} else {
		// Only update input if we're not in chat mode or if we haven't handled the key already
		if m.mode != ChatMode {
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			cmds = append(cmds, inputCmd)
		}
	}

	// Update viewports
	if m.mode == ChatMode {
		m.viewport, _ = m.viewport.Update(msg)
	} else {
		m.filesViewport, _ = m.filesViewport.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleChatKeys(msg tea.KeyMsg, cmds *[]tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "ctrl+f":
		m.mode = FileSelectMode
		m.input.Placeholder = "Space: select, Enter: open folder, Esc: back to chat"
		*cmds = append(*cmds, loadFiles(m.currentDir))
		return m, tea.Batch(*cmds...)
	case "enter":
		text := strings.TrimSpace(m.input.Value())
		if text == "" || m.isWaiting {
			return m, nil
		}

		m.isWaiting = true
		m.input.SetValue("")

		// Create display version (summary) and AI version (full context)
		contextSummary := m.buildFileContextSummary()
		fullContext := m.buildFileContext()

		// Add user message with summary for display
		userDisplayContent := text
		if contextSummary != "" {
			userDisplayContent = fmt.Sprintf("%s\n\n%s", text, contextSummary)
		}

		// Create full content for AI (with actual file contents)
		userAIContent := text
		if fullContext != "" {
			userAIContent = fmt.Sprintf("%s\n\n--- Context from selected files ---\n%s", text, fullContext)
		}

		// Add to history with display version
		m.history = append(m.history,
			ChatMessage{Role: "user", Content: userDisplayContent},
			ChatMessage{Role: "assistant", Content: ""})
		m.viewport.SetContent(renderHistory(m.history))

		// But send the full context to AI
		*cmds = append(*cmds, fetchResponseWithContext(userAIContent, m.history), m.spinner.Tick)
		return m, tea.Batch(*cmds...)
	case "up":
		m.viewport.ScrollUp(1)
	case "down":
		m.viewport.ScrollDown(1)
	default:
		// Let the input handle all other keys (typing)
		if !m.isWaiting {
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			if inputCmd != nil {
				*cmds = append(*cmds, inputCmd)
			}
		}
	}
	return m, nil
}

func (m Model) handleFileKeys(msg tea.KeyMsg, cmds *[]tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = ChatMode
		m.input.Placeholder = "Enter prompt (Ctrl+F for files, Ctrl+C to exit)"
		m.viewport.SetContent(renderHistory(m.history))
		return m, nil
	case "up", "k":
		if m.filesCursor > 0 {
			m.filesCursor--
			m.filesViewport.SetContent(m.renderFiles())
		}
	case "down", "j":
		if m.filesCursor < len(m.files)-1 {
			m.filesCursor++
			m.filesViewport.SetContent(m.renderFiles())
		}
	case " ": // Space to toggle selection
		if len(m.files) > 0 && !m.files[m.filesCursor].IsDir {
			filePath := m.files[m.filesCursor].Path

			// Toggle selection
			found := false
			for i, selected := range m.selectedFiles {
				if selected == filePath {
					m.selectedFiles = append(m.selectedFiles[:i], m.selectedFiles[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				m.selectedFiles = append(m.selectedFiles, filePath)
			}

			m.files[m.filesCursor].Selected = !found
			m.filesViewport.SetContent(m.renderFiles())
		}
	case "enter":
		if len(m.files) > 0 && m.files[m.filesCursor].IsDir {
			m.currentDir = m.files[m.filesCursor].Path
			*cmds = append(*cmds, loadFiles(m.currentDir))
			return m, tea.Batch(*cmds...)
		}
	}
	return m, nil
}

// Build the full context for AI (with file contents)
func (m Model) buildFileContext() string {
	if len(m.selectedFiles) == 0 {
		return ""
	}

	var context strings.Builder
	for _, filePath := range m.selectedFiles {
		content, err := readFileContent(filePath)
		if err != nil {
			context.WriteString(fmt.Sprintf("Error reading %s: %s\n\n", filePath, err.Error()))
			continue
		}

		// Compress/optimize content
		optimizedContent := optimizeFileContent(content, filePath)

		context.WriteString(fmt.Sprintf("File: %s\n", filePath))
		context.WriteString("```\n")
		context.WriteString(optimizedContent)
		context.WriteString("\n```\n\n")
	}

	return context.String()
}

// Build a summary for display in chat history (without file contents)
func (m Model) buildFileContextSummary() string {
	if len(m.selectedFiles) == 0 {
		return ""
	}

	var summary strings.Builder
	summary.WriteString("Attached files:\n")
	for _, filePath := range m.selectedFiles {
		fileName := filepath.Base(filePath)
		// Show relative path if it's in current working directory
		relPath, err := filepath.Rel(".", filePath)
		if err != nil {
			relPath = filePath
		}
		summary.WriteString(fmt.Sprintf("  • %s (%s)\n", fileName, relPath))
	}

	return summary.String()
}

// Optimize file content for AI consumption
func optimizeFileContent(content, filePath string) string {
	// Remove excessive whitespace and empty lines
	lines := strings.Split(content, "\n")
	var optimized []string

	emptyLineCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip excessive empty lines (keep max 1 consecutive empty line)
		if trimmed == "" {
			emptyLineCount++
			if emptyLineCount <= 1 {
				optimized = append(optimized, "")
			}
			continue
		}
		emptyLineCount = 0

		// Remove comments for certain file types to reduce size (optional)
		if shouldStripComments(filePath) {
			if strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "/*") ||
				strings.HasPrefix(trimmed, "*") {
				continue
			}
		}

		optimized = append(optimized, line) // Keep original indentation
	}

	result := strings.Join(optimized, "\n")

	// If still too large, truncate with message
	const maxSize = 8000 // Adjust based on your needs
	if len(result) > maxSize {
		truncated := result[:maxSize-200] // Leave room for message
		result = truncated + "\n\n... (file truncated for brevity) ..."
	}

	return result
}

// Decide whether to strip comments based on file type
func shouldStripComments(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	stripCommentExts := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c"}

	for _, stripExt := range stripCommentExts {
		if ext == stripExt {
			return true
		}
	}
	return false
}

func (m Model) renderFiles() string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("Directory: %s\n", m.currentDir))
	content.WriteString(fmt.Sprintf("Selected files: %d\n\n", len(m.selectedFiles)))

	for i, file := range m.files {
		prefix := "  "
		if i == m.filesCursor {
			prefix = "> "
		}

		var line string
		if file.IsDir {
			line = folderStyle.Render(fmt.Sprintf("%s◯ %s/", prefix, file.Name))
		} else {
			fileIcon := "◯"
			if file.Selected {
				fileIcon = ""
				line = checkedStyle.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			} else {
				line = fileStyle.Render(fmt.Sprintf("%s%s %s", prefix, fileIcon, file.Name))
			}
		}

		if i == m.filesCursor {
			line = selectedStyle.Render(line)
		}

		content.WriteString(line + "\n")
	}

	return content.String()
}

func (m Model) View() string {
	if m.mode == FileSelectMode {
		outputBox := outputStyle.Render(m.filesViewport.View())
		inputBox := inputStyle.Render(m.input.View())
		helpText := helpStyle.Render("↑/↓: navigate, Space: select file, Enter: open folder, Esc: back to chat")
		return fmt.Sprintf("%s\n%s\n%s", outputBox, inputBox, helpText)
	}

	outputBox := outputStyle.Render(m.viewport.View())
	var spinnerLine string
	inputBox := inputStyle.Render(m.input.View())

	var helpText string
	if len(m.selectedFiles) > 0 {
		helpText = helpStyle.Render(fmt.Sprintf("Context: %d files selected | Ctrl+F: file browser", len(m.selectedFiles)))
	} else {
		helpText = helpStyle.Render("Ctrl+F: file browser, ↑/↓: scroll")
	}

	return fmt.Sprintf("%s\n%s%s\n%s", outputBox, spinnerLine, inputBox, helpText)
}
