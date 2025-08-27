package files

type ChatMessage struct {
	Role    string
	Content string
}

type FileItem struct {
	Name     string
	Path     string
	IsDir    bool
	Selected bool
}

type FileChange struct {
	FilePath     string
	OriginalCode string
	ProposedCode string
	Description  string
	Accepted     *bool
	LineStart    int
	LineEnd      int
}
