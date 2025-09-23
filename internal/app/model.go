package app

import (
	"fmt"
	"getctx/internal/logger"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	config                *Config
	fsys                  FileSystem
	pathInput             textinput.Model
	isInputMode           bool
	isFilterMode          bool
	filterQuery           string
	inputErrorMsg         string
	viewport              viewport.Model
	width                 int
	height                int
	completionSuggestions []string
}

func NewModel(startPath string, config *Config, fsys FileSystem) (*Model, error) {
	path, err := fsys.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for '%s': %w", startPath, err)
	}

	items, err := loadItems(fsys, path, config)
	if err != nil {
		return nil, fmt.Errorf("could not read directory '%s': %w", path, err)
	}

	ti := textinput.New()
	ti.Prompt = Icons.Cursor + Elements.List.CursorEmpty
	ti.Focus()

	vp := viewport.New(0, 0)

	m := &Model{
		path:        path,
		items:       items,
		selected:    make(map[string]struct{}),
		config:      config,
		fsys:        fsys,
		pathInput:   ti,
		isInputMode: false,
		viewport:    vp,
	}

	m.viewport.SetContent(m.renderFileList())
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputWidth := max(m.width-len(m.pathInput.Prompt)-1, 1)
		m.pathInput.Width = inputWidth
	}

	if m.isInputMode {
		// --- POCZĄTEK KOREKTY LOGIKI ---
		oldValue := m.pathInput.Value()

		// Najpierw obsługujemy nasze specjalne skróty klawiszowe (TAB, Enter, ESC).
		// Funkcja zwróci `true`, jeśli klawisz został obsłużony.
		cmd, keyWasHandled := m.handleInputModeKeys(msg)
		cmds = append(cmds, cmd)

		// Jeśli klawisz NIE był specjalnym skrótem, przekazujemy go do pola tekstowego,
		// aby mogło ono dodać literę do swojej wartości.
		if !keyWasHandled {
			m.pathInput, cmd = m.pathInput.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Po wszystkich operacjach sprawdzamy, czy tekst w polu się zmienił.
		if m.pathInput.Value() != oldValue {
			// Jeśli tak, natychmiast odświeżamy listę sugestii.
			m.updateCompletions()
		}
		// --- KONIEC KOREKTY LOGIKI ---

	} else if m.isFilterMode {
		m.pathInput, cmd = m.pathInput.Update(msg)
		cmds = append(cmds, cmd)
		cmd = m.updateFilterMode(msg)
	} else {
		cmd = m.updateNormalMode(msg)
		cmds = append(cmds, cmd)
	}

	// Reszta funkcji pozostaje bez zmian
	headerContent := m.renderHeader()
	footerContent := m.renderFooter()
	headerHeight := lipgloss.Height(headerContent)
	footerHeight := lipgloss.Height(footerContent)
	m.viewport.Width = m.width
	m.viewport.Height = m.height - headerHeight - footerHeight
	m.viewport.SetContent(m.renderFileList())
	m.ensureCursorVisible()

	return m, tea.Batch(cmds...)
}

func (m *Model) changeDirectory(newPath string) {
	newItems, err := loadItems(m.fsys, newPath, m.config)
	if err != nil {
		logger.Error("changeDirectory.loadItems", map[string]any{
			"message": "Failed to load directory items",
			"path":    newPath,
			"error":   err.Error(),
		})
		return
	}

	m.path = newPath
	m.items = newItems
	m.cursor = 0
	m.filterQuery = ""
	m.pathInput.Reset()
	m.viewport.GotoTop()
}

func loadItems(fsys FileSystem, path string, config *Config) ([]item, error) {
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
