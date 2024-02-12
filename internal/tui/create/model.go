package create

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitlab.com/l0nax/changelog-go/internal/config"
)

type Entry struct {
	Title string
	Type  config.ChangeType
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// entryItem is wrapper, implementing [list.Item] of [config.ChangeType].
type entryItem struct {
	ct config.ChangeType
}

func (i entryItem) Title() string       { return i.ct.Title }
func (i entryItem) Description() string { return i.ct.Description.UnwrapOrZero() }
func (i entryItem) FilterValue() string { return i.ct.Title + i.Description() }

// model is the TUI model of the create interface.
type model struct {
	// entryList is the list of the available change type entries.
	entryList list.Model

	// typeSelected holds the state information whether the user selected
	// a change type.
	typeSelected bool
	selectedType config.ChangeType
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if msg.Type == tea.KeyEnter {
			return m.handleEnter()
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.entryList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd

	// only update if we're in the correct stage.
	if !m.typeSelected {
		m.entryList, cmd = m.entryList.Update(msg)
	}

	return m, cmd
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch {
	case !m.typeSelected:
		m.typeSelected = true

		entry := m.entryList.SelectedItem().(entryItem)
		m.selectedType = entry.ct

		return m, tea.Quit
	}

	return m, nil
}

func (m model) generateView() string {
	var buf strings.Builder
	buf.Grow(2048)

	if !m.typeSelected {
		panic("never call generateView if selectedType is not set")
	}

	return buf.String()
}

func (m model) View() string {
	switch {
	case !m.typeSelected:
		// let the user select the type
		return docStyle.Render(m.entryList.View())
	}

	return docStyle.Render(m.entryList.View())
}

func Run() (Entry, error) {
	// NOTE: We need to clone to not change config.C.Entry.Types
	rawItems := slices.Clone(config.C.Entry.Types)
	rawItems = slices.DeleteFunc(rawItems, func(tt config.ChangeType) bool {
		return tt.Hidden
	})
	items := make([]list.Item, len(rawItems))

	for i := range rawItems {
		items[i] = entryItem{
			ct: rawItems[i],
		}
	}

	m := model{
		entryList: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.entryList.Title = "Select a change type"

	p := tea.NewProgram(m, tea.WithAltScreen())

	end, err := p.Run()
	if err != nil {
		return Entry{}, err
	}

	endModel := end.(model)

	return Entry{
		Type: endModel.selectedType,
	}, nil
}

var _ list.Item = (*entryItem)(nil)
