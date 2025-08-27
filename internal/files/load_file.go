package files

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func LoadFiles(dir string) tea.Cmd {
	return func() tea.Msg {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		var fileItems []FileItem
		if dir != "/" && dir != "." {
			fileItems = append(fileItems, FileItem{
				Name:  "..",
				Path:  filepath.Dir(dir),
				IsDir: true,
			})
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			fileItems = append(fileItems, FileItem{
				Name:  file.Name(),
				Path:  filepath.Join(dir, file.Name()),
				IsDir: file.IsDir(),
			})
		}
		return fileItems
	}
}

func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
