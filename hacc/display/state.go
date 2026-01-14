package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Event interface{}

type UpEvent struct{}
type DownEvent struct{}
type LeftEvent struct{}
type RightEvent struct{}
type EnterEvent struct {
	ServiceName string
}
type BackEvent struct{}
type NumberEvent struct {
	Number int
}
type RuneEvent struct{}

// //////////////////////////////////////////////////////////////////////////
// Bubbletea Model main Update method
// //////////////////////////////////////////////////////////////////////////
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			return m.state.Update(m, UpEvent{})
		case tea.KeyDown:
			return m.state.Update(m, DownEvent{})
		case tea.KeyLeft:
			return m.state.Update(m, LeftEvent{})
		case tea.KeyRight:
			return m.state.Update(m, RightEvent{})
		case tea.KeyEnter:
			return m.state.Update(m, EnterEvent{})
		case tea.KeyBackspace:
			return m.state.Update(m, BackEvent{})
		case tea.KeyRunes:
			// Check if single digit number for List selection
			if len(msg.Runes) == 1 && msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
				n := int(msg.Runes[0] - '0')
				return m.state.Update(m, NumberEvent{Number: n})
			}
			return m.state.Update(m, RuneEvent{})
		}
	}
	return m, nil
}

// ///////////////////////////////////////////////////////////////////////////
// State interface and common functions for UI
// //////////////////////////////////////////////////////////////////////
type State interface {
	Update(m model, e Event) (model, tea.Cmd)
}
type Backable interface {
	Back(m model) (model, tea.Cmd)
}

type EnterFunc func(m model) (model, tea.Cmd)

func updateListState(m model, e Event, numItems int, onEnter EnterFunc) (model, tea.Cmd) {
	switch e.(type) {
	case UpEvent:
		return CursorUp(m, numItems), nil
	case DownEvent:
		return CursorDown(m, numItems), nil
	case LeftEvent:
		return CursorLeft(m, numItems), nil
	case RightEvent:
		return CursorRight(m, numItems), nil
	case EnterEvent:
		return onEnter(m)
	case BackEvent:
		return m.state.(Backable).Back(m)
	case NumberEvent:
		return CursorNumber(m, e.(NumberEvent).Number), nil
	default:
		return m, nil
	}
}

func CursorUp(m model, numItems int) model {
	if m.cursor > 0 {
		m.cursor--
		return m
	}
	// Consider less-than-full last page for wrapping
	if m.page == NumPages(numItems, m.pageSize)-1 {
		m.cursor = NumLastPageItems(numItems, m.pageSize) - 1
		return m
	}
	m.cursor = m.pageSize - 1
	return m
}

func CursorDown(m model, numItems int) model {
	if m.page == NumPages(numItems, m.pageSize)-1 {
		// Consider less-than-full last page for wrapping
		if m.cursor < NumLastPageItems(numItems, m.pageSize)-1 {
			m.cursor++
			return m
		}
		m.cursor = 0
		return m
	}
	// Full page case
	if m.cursor < m.pageSize-1 {
		m.cursor++
		return m
	}
	m.cursor = 0
	return m
}

func CursorNumber(m model, n int) model {
	if n < 0 || n > m.pageSize {
		return m
	}
	// Consider less-than-full last page
	if m.page == NumPages(len(m.vault.Services), m.pageSize)-1 {
		if n >= NumLastPageItems(len(m.vault.Services), m.pageSize) {
			return m
		}
	}
	m.cursor = n - 1
	return m
}

func CursorLeft(m model, numItems int) model {
	if m.page > 0 {
		m.page--
		return m
	} else {
		// Wrap around to last page and adjust cursor if needed
		m.page = NumPages(numItems, m.pageSize) - 1
		if m.cursor >= NumLastPageItems(numItems, m.pageSize) {
			m.cursor = NumLastPageItems(numItems, m.pageSize) - 1
		}
	}
	return m
}

func CursorRight(m model, numItems int) model {
	if m.page < NumPages(numItems, m.pageSize)-1 {
		m.page++
	} else {
		// wrap around to first page
		m.page = 0
	}
	if m.page == NumPages(numItems, m.pageSize)-1 {
		// Adjust cursor for less-than-full last page
		if m.cursor >= NumLastPageItems(numItems, m.pageSize) {
			m.cursor = NumLastPageItems(numItems, m.pageSize) - 1
		}
	}
	return m
}

func NumPages(totalItems, pageSize int) int {
	if totalItems%pageSize == 0 {
		return totalItems / pageSize
	}
	return totalItems/pageSize + 1
}

func NumLastPageItems(totalItems, pageSize int) int {
	if totalItems%pageSize == 0 {
		return pageSize
	}
	return totalItems % pageSize
}

// ///////////////////////////////////////////////////////////////////////////
// Concrete State structs and Update methods
// ////////////////////////////////////////////////////////////////////////
type WelcomeState struct{}

func (s WelcomeState) Update(m model, e Event) (model, tea.Cmd) {
	if len(m.vault.Services) == 0 {
		m.message = "No credentials are stored in the Vault yet! Add one to get started."
		m.state = &EndState{}
		return m, nil
	}
	// Any key event moves to service list
	m.state = &ServiceListState{}
	return m, nil
}

// Pressing back from welcome state quits the program
func (s WelcomeState) Back(m model) (model, tea.Cmd) {
	return m, tea.Quit
}

type CredentialState struct{}

func (s CredentialState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case BackEvent:
		return s.Back(m)
	}
	return m, nil
}

func (s CredentialState) Back(m model) (model, tea.Cmd) {
	m.state = &UsernameListState{}
	return m, nil
}

type MessageState struct{}

func (s MessageState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case BackEvent:
		return s.Back(m)
	}
	return m, nil
}

func (s MessageState) Back(m model) (model, tea.Cmd) {
	m.state = &ServiceListState{}
	return m, nil
}

type EndState struct{}

func (s EndState) Update(m model, e Event) (model, tea.Cmd) {
	return m, tea.Quit
}

type ServiceListState struct{}

func (s ServiceListState) Update(m model, e Event) (model, tea.Cmd) {
	return updateListState(
		m,
		e,
		len(m.vault.Services),
		func(m model) (model, tea.Cmd) {
			if len(m.vault.Services) == 0 {
				m.message = "No services available to copy."
				return m, nil
			}
			// Get the selected service name
			services := m.vault.ListServices(m.serviceName)
			selectedService := services[m.page*m.pageSize+m.cursor]
			m.serviceName = selectedService
			m.cursor = 0
			m.state = &UsernameListState{}
			return m, nil
		},
	)
}

func (s ServiceListState) Back(m model) (model, tea.Cmd) {
	return m, tea.Quit
}

type UsernameListState struct{}

func (s UsernameListState) Update(m model, e Event) (model, tea.Cmd) {
	return updateListState(
		m,
		e,
		len(m.vault.Services[m.serviceName].GetUsers("")),
		func(m model) (model, tea.Cmd) {
			if len(m.vault.Services) == 0 {
				m.message = "No services available to copy."
				return m, nil
			}
			// Get the selected service name
			service := m.vault.Services[m.serviceName]
			users := service.GetUsers("")
			selectedUser := users[m.page*m.pageSize+m.cursor]
			m.username = selectedUser
			m.state = &CredentialState{}
			return m, nil
		},
	)
}

func (s UsernameListState) Back(m model) (model, tea.Cmd) {
	m.state = &ServiceListState{}
	m.serviceName = ""
	return m, nil
}
