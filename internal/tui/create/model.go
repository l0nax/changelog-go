package create

import (
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitlab.com/l0nax/changelog-go/internal/config"
	"gitlab.com/l0nax/changelog-go/internal/tui/common"
)

// ErrCanceled is returned when the user canceled the operation.
var ErrCanceled = errors.New("canceled")

const (
	DefaultValidateOkPrefix  = "✔"
	DefaultValidateErrPrefix = "✘"

	ansiColorValidateOK  = "2" // ansiColorValidateOK is the OK ANSI color (green)
	ansiColorValidateErr = "1" // ansiColorValidateErr is the Error ANSI color (red)
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

	canceled bool

	// typeSelected holds the state information whether the user selected
	// a change type.
	typeSelected bool
	selectedType config.ChangeType

	askTitle       bool
	title          string
	titleInput     textinput.Model
	showTitleError bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.canceled = true

			return m, tea.Quit

		case tea.KeyEnter:
			return m.handleEnter(msg)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.entryList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd

	// only update if we're in the correct stage.
	if !m.typeSelected {
		m.entryList, cmd = m.entryList.Update(msg)
	} else if m.askTitle {
		m.titleInput, cmd = m.titleInput.Update(msg)
		m.title = m.titleInput.Value()

		// reset error once user enter something
		m.showTitleError = m.title == ""
	}

	return m, cmd
}

func (m model) handleEnter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case !m.typeSelected:
		m.typeSelected = true

		entry := m.entryList.SelectedItem().(entryItem)
		m.selectedType = entry.ct

		slog.Debug("Selected change type", slog.Any("change_type", m.selectedType))

		// if we have to ask for the title, proceed
		if m.askTitle && m.title == "" {
			return m, textinput.Blink
		}

		return m, tea.Quit

	case m.askTitle && msg.Type == tea.KeyEnter:
		// only allow exit if it is non-empty
		if m.title != "" {
			return m, tea.Quit
		}

		m.showTitleError = true
	}

	return m, nil
}

func (m model) View() string {
	switch {
	case !m.typeSelected:
		// let the user select the type
		return docStyle.Render(m.entryList.View())

	case m.askTitle:
		var buf strings.Builder
		buf.Grow(1024)

		if m.showTitleError {
			buf.WriteString(common.FontColor(DefaultValidateErrPrefix, ansiColorValidateErr))
		} else {
			buf.WriteString(common.FontColor(DefaultValidateOkPrefix, ansiColorValidateOK))
		}

		buf.WriteByte(' ')

		buf.WriteString("Please enter the changelog entry title\n\n")

		buf.WriteString(m.titleInput.View())
		buf.WriteString("\n\n")
		buf.WriteString("(press enter to continue)\n\n")

		if m.showTitleError {
			buf.WriteString(common.FontColor("Please enter a changelog title!\n", ansiColorValidateErr))
		}

		return buf.String()
	}

	return docStyle.Render(m.entryList.View())
}

// Canceled determine whether the operation is cancelled
func (m model) Canceled() bool {
	return m.canceled
}

// Run shows the TUI and returns the selected change type.
// If askTitle is set to true, the user will be asked to enter
// a non-empty title string.
//
// The [ErrCanceled] is returned if the user canceled, i.e. hits Ctrl-C.
func Run(askTitle bool) (Entry, error) {
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
		askTitle:  askTitle,
	}
	m.entryList.Title = "Please select a change type"

	if askTitle {
		m.titleInput = textinput.New()
		m.titleInput.Placeholder = "Title"
		m.titleInput.Focus()

		m.showTitleError = true
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	end, err := p.Run()
	if err != nil {
		return Entry{}, err
	}

	endModel := end.(model)
	if endModel.Canceled() {
		return Entry{}, ErrCanceled
	}

	return Entry{
		Type:  endModel.selectedType,
		Title: endModel.title,
	}, nil
}

var _ list.Item = (*entryItem)(nil)
