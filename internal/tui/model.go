package tui

import (
	"fmt"
	"getctx/internal/config"
	"getctx/internal/fs"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	name       string
	isDir      bool
	isExcluded bool
}

type Model struct {
	path                  string
	items                 []item
	cursor                int
	selected              map[string]struct{}
	config                *config.Config
	fsys                  fs.FileSystem
	pathInput             textinput.Model
	isInputMode           bool
	isFilterMode          bool
	filterQuery           string
	inputErrorMsg         string
	viewport              viewport.Model
	completionViewport    viewport.Model
	width                 int
	height                int
	completionSuggestions []string
	Aborted               bool
}

func NewModel(startPath string, config *config.Config, fsys fs.FileSystem) (*Model, error) {
	path, err := fsys.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", startPath, err)
	}

	items, err := loadItems(fsys, path, config)
	if err != nil {
		return nil, fmt.Errorf("could not read directory '%s': %w", path, err)
	}

	ti := textinput.New()
	ti.Prompt = Icons.Cursor + " "
	ti.Focus()

	vp := viewport.New(0, 0)
	completionVP := viewport.New(0, 0)

	m := &Model{
		path:               path,
		items:              items,
		selected:           make(map[string]struct{}),
		config:             config,
		fsys:               fsys,
		pathInput:          ti,
		viewport:           vp,
		completionViewport: completionVP,
		Aborted:            false,
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) GetSelectedPaths() []string {
	paths := make([]string, 0, len(m.selected))
	for path := range m.selected {
		paths = append(paths, path)
	}
	return paths
}

func (m *Model) changeDirectory(newPath string) {
	newItems, err := loadItems(m.fsys, newPath, m.config)
	if err != nil {
		m.inputErrorMsg = "Error reading directory: " + err.Error()
		return
	}

	m.path = newPath
	m.items = newItems
	m.cursor = 0
	m.filterQuery = ""
	m.pathInput.Reset()
	m.viewport.GotoTop()
}

func loadItems(fsys fs.FileSystem, path string, config *config.Config) ([]item, error) {
	dirEntries, err := fsys.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var items []item
	for _, entry := range dirEntries {
		items = append(items, item{
			name:       entry.Name(),
			isDir:      entry.IsDir(),
			isExcluded: config.IsExcluded(entry.Name()),
		})
	}
	return items, nil
}

func (m *Model) clampCursor() {
	visibleItems := m.getVisibleItems()
	maxCursor := max(len(visibleItems)-1, 0)
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}
}
