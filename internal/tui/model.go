package tui

import (
	"fmt"

	"github.com/kacperzielinskidev/getctx/internal/config"
	"github.com/kacperzielinskidev/getctx/internal/fs"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiMode int

const (
	modeNormal tuiMode = iota
	modePathInput
	modeFilter
)

type listItem struct {
	name       string
	isDir      bool
	isExcluded bool
}

type Model struct {
	config                *config.Config
	fsys                  fs.FileSystem
	path                  string
	items                 []listItem
	selected              map[string]struct{}
	Aborted               bool
	textInput             textinput.Model
	viewport              viewport.Model
	completionViewport    viewport.Model
	mode                  tuiMode
	cursor                int
	filterQuery           string
	inputErrorMsg         string
	completionSuggestions []string
	width                 int
	height                int
}

func NewModel(startPath string, config *config.Config, fsys fs.FileSystem) (*Model, error) {
	path, err := fsys.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", startPath, err)
	}

	items, err := loadListItems(fsys, path, config)
	if err != nil {
		return nil, fmt.Errorf("could not read directory '%s': %w", path, err)
	}

	ti := textinput.New()
	ti.Prompt = Icons.Cursor + Elements.List.CursorEmpty
	ti.Focus()

	m := &Model{
		config:             config,
		fsys:               fsys,
		path:               path,
		items:              items,
		selected:           make(map[string]struct{}),
		textInput:          ti,
		viewport:           viewport.New(0, 0),
		completionViewport: viewport.New(0, 0),
		mode:               modeNormal,
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
	newItems, err := loadListItems(m.fsys, newPath, m.config)
	if err != nil {
		m.inputErrorMsg = "Error reading directory: " + err.Error()
		return
	}

	m.path = newPath
	m.items = newItems
	m.cursor = 0
	m.filterQuery = ""
	m.textInput.Reset()
	m.viewport.GotoTop()
}

func loadListItems(fsys fs.FileSystem, path string, config *config.Config) ([]listItem, error) {
	dirEntries, err := fsys.ReadDir(path)
	if err != nil {
		return nil, err
	}

	items := make([]listItem, len(dirEntries))
	for i, entry := range dirEntries {
		items[i] = listItem{
			name:       entry.Name(),
			isDir:      entry.IsDir(),
			isExcluded: config.IsExcluded(entry.Name()),
		}
	}
	return items, nil
}

func (m *Model) clampCursor() {
	visibleItems := m.getVisibleItems()
	maxCursor := len(visibleItems) - 1
	maxCursor = max(maxCursor, 0)

	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	m.cursor = max(m.cursor, 0)
}
