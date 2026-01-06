package main

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

// ///////////////////////////////////////////////////////////////////////////
// State interface and common functions for UI
// //////////////////////////////////////////////////////////////////////
type State interface {
	Update(m model, e Event) (model, tea.Cmd)
}

func CursorUp(m model) model {
	if m.cursor > 0 {
		m.cursor--
		return m
	}
	// Consider less-than-full last page for wrapping
	if m.page == NumPages(len(m.vault.services), m.pageSize)-1 {
		m.cursor = NumLastPageItems(len(m.vault.services), m.pageSize) - 1
		return m
	}
	m.cursor = m.pageSize - 1
	return m
}

func CursorDown(m model) model {
	if m.page == NumPages(len(m.vault.services), m.pageSize)-1 {
		// Consider less-than-full last page for wrapping
		if m.cursor < NumLastPageItems(len(m.vault.services), m.pageSize)-1 {
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
	if m.page == NumPages(len(m.vault.services), m.pageSize)-1 {
		if n >= NumLastPageItems(len(m.vault.services), m.pageSize) {
			return m
		}
	}
	m.cursor = n - 1
	return m
}

func Back(m model) model {
	m.state = &ServiceListState{}
	m.serviceName = ""
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
	if len(m.vault.services) == 0 {
		m.message = "No credentials are stored in the Vault yet! Add one to get started."
		m.state = &EndState{}
		return m, nil
	}
	m.state = &ServiceListState{}
	return m, nil
}

type CredentialState struct{}

func (s CredentialState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case BackEvent:
		return Back(m), nil
	}
	return m, nil
}

type MessageState struct{}

func (s MessageState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case BackEvent:
		return Back(m), nil
	}
	return m, nil
}

type EndState struct{}

func (s EndState) Update(m model, e Event) (model, tea.Cmd) {
	return m, tea.Quit
}

type ServiceListState struct{}

func (s ServiceListState) Update(m model, e Event) (model, tea.Cmd) {
	switch e.(type) {
	case UpEvent:
		return CursorUp(m), nil
	case DownEvent:
		return CursorDown(m), nil
	case LeftEvent:
		if m.page > 0 {
			m.page--
			return m, nil
		} else {
			// Wrap around to last page and adjust cursor if needed
			m.page = NumPages(len(m.vault.services), m.pageSize) - 1
			if m.cursor >= NumLastPageItems(len(m.vault.services), m.pageSize) {
				m.cursor = NumLastPageItems(len(m.vault.services), m.pageSize) - 1
			}
		}
		return m, nil
	case RightEvent:
		if m.page < NumPages(len(m.vault.services), m.pageSize)-1 {
			m.page++
		} else {
			// wrap around to first page
			m.page = 0
		}
		if m.page == NumPages(len(m.vault.services), m.pageSize)-1 {
			// Adjust cursor for less-than-full last page
			if m.cursor >= NumLastPageItems(len(m.vault.services), m.pageSize) {
				m.cursor = NumLastPageItems(len(m.vault.services), m.pageSize) - 1
			}
		}
		return m, nil
	case NumberEvent:
		return CursorNumber(m, e.(NumberEvent).Number), nil
	case EnterEvent:
		if len(m.vault.services) == 0 {
			m.message = "No services available to copy."
			return m, nil
		}
		// Get the selected service name
		services := m.vault.ListServices(m.serviceName)
		selectedService := services[m.page*m.pageSize+m.cursor]
		m.serviceName = selectedService
		m.state = &UsernameListState{}
		return m, nil

	}
	return m, nil
}

type UsernameListState struct{}

func (s UsernameListState) Update(m model, e Event) (model, tea.Cmd) {
	service := m.vault.services[m.serviceName]
	switch e.(type) {
	case BackEvent:
		return Back(m), nil
	case UpEvent:
		return CursorUp(m), nil
	case DownEvent:
		return CursorDown(m), nil
	case LeftEvent:
		if m.page > 0 {
			m.page--
			return m, nil
		} else {
			// Wrap around to last page and adjust cursor if needed
			m.page = NumPages(len(service.GetUsers("")), m.pageSize) - 1
			if m.cursor >= NumLastPageItems(len(service.GetUsers("")), m.pageSize) {
				m.cursor = NumLastPageItems(len(service.GetUsers("")), m.pageSize) - 1
			}
		}
		return m, nil
	case RightEvent:
		if m.page < NumPages(len(service.GetUsers("")), m.pageSize)-1 {
			m.page++
		} else {
			// wrap around to first page
			m.page = 0
		}
		if m.page == NumPages(len(service.GetUsers("")), m.pageSize)-1 {
			// Adjust cursor for less-than-full last page
			if m.cursor >= NumLastPageItems(len(service.GetUsers("")), m.pageSize) {
				m.cursor = NumLastPageItems(len(service.GetUsers("")), m.pageSize) - 1
			}
		}
		return m, nil
	case NumberEvent:
		return CursorNumber(m, e.(NumberEvent).Number), nil
	case EnterEvent:
		if len(m.vault.services) == 0 {
			m.message = "No services available to copy."
			return m, nil
		}
		// Get the selected service name
		users := service.GetUsers("")
		selectedUser := users[m.page*m.pageSize+m.cursor]
		m.username = selectedUser
		m.state = &CredentialState{}
		return m, nil
	}
	return m, nil
}
