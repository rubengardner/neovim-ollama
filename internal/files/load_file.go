package files

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rubengardner/neovim-ollama/internal/model"
)

func loadFiles(dir string) tea.Cmd {
	return func() tea.Msg {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		var fileItems []model.FileItem

		// Add parent directory if not at root
		if dir != "/" && dir != "." {
			fileItems = append(fileItems, model.FileItem{
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

			fileItems = append(fileItems, model.FileItem{
				Name:  file.Name(),
				Path:  filepath.Join(dir, file.Name()),
				IsDir: file.IsDir(),
			})
		}

		return model.FilesLoadedMsg(fileItems)
	}
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
