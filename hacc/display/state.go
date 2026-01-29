package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nbailey20/hacc/hacc/cli"
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
type SpaceEvent struct{}
type NumberEvent struct {
	Number int
}
type RuneEvent struct {
	Rune rune
}

// //////////////////////////////////////////////////////////////////////////
// Bubbletea Model main Update method
// //////////////////////////////////////////////////////////////////////////
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			return m.state.Update(m, UpEvent{})
		case "down":
			return m.state.Update(m, DownEvent{})
		case "left":
			return m.state.Update(m, LeftEvent{})
		case "right":
			return m.state.Update(m, RightEvent{})
		case "enter":
			return m.state.Update(m, EnterEvent{})
		case "backspace":
			return m.state.Update(m, BackEvent{})
		case " ":
			return m.state.Update(m, SpaceEvent{})
		default:
			if len(msg.String()) != 1 {
				return m, nil
			}
			// Check if single digit number for List selection
			key := msg.String()[0]
			if key >= '0' && key <= '9' {
				n := int(key - '0')
				return m.state.Update(m, NumberEvent{Number: n})
			}
			return m.state.Update(m, RuneEvent{Rune: rune(key)})
		}
	case ErrorMsg:
		m.endSuccess = false
		m.endError = msg.ErrorValue()
		m.endMessage = msg.DisplayValue()
	case SuccessMsg:
		m.endSuccess = true
		m.endError = nil
		m.endMessage = msg.DisplayValue()
	case PasswordGeneratedMsg:
		m.password = msg.Password
	case PasswordLoadedMsg:
		m.password = msg.Password
		return m, copyPasswordCmd(m.password)
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
	switch e := e.(type) {
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
		return CursorNumber(m, e.Number), nil
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
	if n < 1 || n > m.pageSize {
		return m
	}
	numResources := len(m.vault.Services)
	if _, ok := m.state.(*UsernameListState); ok {
		numResources = m.vault.Services[m.serviceName].NumUsers()
	}
	// Consider less-than-full last page
	if m.page == NumPages(numResources, m.pageSize)-1 {
		if n > NumLastPageItems(numResources, m.pageSize) {
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
		m.state = &EmptyState{}
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

type EndState struct{}

func (s EndState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case BackEvent:
		return s.Back(m)
	}
	return m, nil
}

func (s EndState) Back(m model) (model, tea.Cmd) {
	m.action = cli.SearchAction{}
	m.state = &UsernameListState{}
	// if the last user in a service was just deleted,
	// update model and go back to the other services
	if !m.vault.HasService(m.serviceName) {
		m.serviceName = ""
		m.username = ""
		m.password = ""
		m.state = &ServiceListState{}
	}
	return m, nil
}

type EmptyState struct{}

func (s EmptyState) Update(m model, e Event) (model, tea.Cmd) {
	return m, tea.Quit
}

type CredentialState struct{}

func (s CredentialState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case SpaceEvent:
		m.showPass = !m.showPass
	case BackEvent:
		return s.Back(m)
	}
	return m, nil
}

func (s CredentialState) Back(m model) (model, tea.Cmd) {
	m.state = &UsernameListState{}
	return m, nil
}

type ConfirmState struct{}

func (s ConfirmState) Update(m model, e Event) (model, tea.Cmd) {
	switch e := e.(type) {
	case EnterEvent:
		m.state = &EndState{}
		return m, addCredentialCmd(m.vault, m.serviceName, m.username, m.password)
	case RuneEvent:
		if e.Rune == 'y' {
			m.state = &EndState{}
			return m, addCredentialCmd(m.vault, m.serviceName, m.username, m.password)
		}
		if e.Rune == 'n' {
			return m, generatePasswordCmd(
				m.digitsInPass,
				m.specialsInPass,
				m.minPassLen,
				m.maxPassLen,
			)
		}
		return m, nil
	case PasswordGeneratedMsg:
		m.password = e.Password
		return m, nil
	default:
		return m, nil
	}
}

type ServiceListState struct{}

func (s ServiceListState) Update(m model, e Event) (model, tea.Cmd) {
	return updateListState(
		m,
		e,
		len(m.vault.Services),
		func(m model) (model, tea.Cmd) {
			if len(m.vault.Services) == 0 {
				m.state = &EmptyState{}
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
			// Get the selected service name
			service := m.vault.Services[m.serviceName]
			users := service.GetUsers("")
			selectedUser := users[m.page*m.pageSize+m.cursor]
			m.username = selectedUser
			m.state = &CredentialState{}
			return m, loadPasswordCmd(m.vault, m.serviceName, m.username)
		},
	)
}

func (s UsernameListState) Back(m model) (model, tea.Cmd) {
	m.state = &ServiceListState{}
	m.serviceName = ""
	return m, nil
}
