package create

import (
	"slices"

	"gitlab.com/l0nax/changelog-go/internal/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Entry struct {
	Title string
	Type  config.ChangeType
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	ct config.ChangeType
}

func (i item) Title() string       { return i.ct.Title }
func (i item) Description() string { return i.ct.Description.UnwrapOrZero() }
func (i item) FilterValue() string { return i.ct.Title + i.Description() }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run() (Entry, error) {
	// NOTE: We need to clone to not change config.C.Entry.Types
	rawItems := slices.Clone(config.C.Entry.Types)
	rawItems = slices.DeleteFunc(rawItems, func(tt config.ChangeType) bool {
		return tt.Hidden
	})
	items := make([]item, len(rawItems))

	for i := range rawItems {
		items[i] = item{
			ct: rawItems[i],
		}
	}

	m := model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Select a change type"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Run(); err != nil {
		return Entry{}, err
	}

	// TODO: return real data
	return Entry{}, nil
}
