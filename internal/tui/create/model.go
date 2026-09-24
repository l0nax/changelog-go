// Package create implements the terminal interface for creating an entry.
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
	"go.l0nax.org/typact"

	"gitlab.com/l0nax/changelog-go/v2/internal/config"
	"gitlab.com/l0nax/changelog-go/v2/internal/tui/common"
)

// ErrCanceled is returned when the user canceled the operation.
var ErrCanceled = errors.New("canceled")

// The prefixes marking the validation state of the title input.
const (
	DefaultValidateOkPrefix  = "✔"
	DefaultValidateErrPrefix = "✘"

	ansiColorValidateOK  = "2" // ansiColorValidateOK is the OK ANSI color (green)
	ansiColorValidateErr = "1" // ansiColorValidateErr is the Error ANSI color (red)
)

// Entry is the change type and title the user selected.
type Entry struct {
	Title string
	Type  config.ChangeType
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// entryItem adapts a [config.ChangeType] to [list.Item].
type entryItem struct {
	ct config.ChangeType
}

func (i entryItem) Title() string       { return i.ct.Title }
func (i entryItem) Description() string { return i.ct.Description.UnwrapOrZero() }
func (i entryItem) FilterValue() string { return i.ct.Title + i.Description() }

// model is the state of the create interface.
type model struct {
	// entryList is the list of the available change type entries.
	entryList list.Model

	canceled bool

	// typeSelected reports whether the user selected a change type.
	typeSelected bool
	selectedType config.ChangeType

	askTitle       bool
	title          string
	titleInput     textinput.Model
	showTitleError bool
}

func (m model) Init() tea.Cmd {
	// Starting straight on the title prompt means the cursor has to blink
	// without a type selection having happened first.
	if m.typeSelected && m.askTitle {
		return textinput.Blink
	}

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

	if !m.typeSelected {
		m.entryList, cmd = m.entryList.Update(msg)
	} else if m.askTitle {
		m.titleInput, cmd = m.titleInput.Update(msg)
		m.title = m.titleInput.Value()

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

		if m.askTitle && m.title == "" {
			return m, textinput.Blink
		}

		return m, tea.Quit

	case m.askTitle && msg.Type == tea.KeyEnter:
		// a title is required before the user can leave this stage
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

// Canceled reports whether the user canceled the operation.
func (m model) Canceled() bool {
	return m.canceled
}

// Options configures what the TUI asks for.
type Options struct {
	// PreselectedType skips the type selection, which is what happens when
	// the caller already passed --type.
	PreselectedType typact.Option[config.ChangeType]

	// AskTitle prompts for a non-empty title.
	AskTitle bool
}

// Run shows the TUI and returns the change type the user selected. Only what
// opts leaves open is prompted for.
//
// It returns [ErrCanceled] if the user canceled.
func Run(cfg config.Config, opts Options) (Entry, error) {
	askTitle := opts.AskTitle

	// cloned because the filtering below would otherwise alter the config
	rawItems := slices.Clone(cfg.Entry.Types)
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

	if changeType, ok := opts.PreselectedType.Deconstruct(); ok {
		m.typeSelected = true
		m.selectedType = changeType
	}

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
