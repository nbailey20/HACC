package main

import (
	"fmt"
	"sort"
	"sync"
)

// //////////////////////////////////////////////////////////////////////
// state struct and methods
// //////////////////////////////////////////////////////////////////////
type state struct {
	Name   string
	Value  string
	Path   string
	loaded bool // is the local copy up-to-date with the current value
	saved  bool // is the current value saved to the backend
	client *ssmClient
}

// If Value is empty, use existing value from backend
// Otherwise immediately save provided value and use it
func NewState(name string, value string, path string) (state, error) {
	state := state{
		Name:   name,
		Value:  value,
		Path:   path,
		loaded: false,
		saved:  false,
		client: NewSsmClient(),
	}
	if value != "" {
		state.loaded = true
		err := state.Save()
		if err == nil {
			state.saved = true
		}
		return state, err
	}
	return state, nil
}

func (s *state) GetName() string {
	return s.Name
}

// Lazy backend fetching, don't pull the value before its needed
func (s *state) GetValue() string {
	if !s.loaded {
		s.Load()
	}
	return s.Value
}

func (s *state) SetValue(value string) {
	s.Value = value
	s.loaded = true
	s.saved = false
}

func (s *state) Save() error {
	if s.saved {
		return nil
	}
	err := s.client.PutParameter(s.Path+s.Name, s.Value)
	if err == nil {
		s.saved = true
	}
	return err
}

func (s *state) Load() error {
	if s.loaded {
		return nil
	}
	value, err := s.client.GetParameter(s.Path + s.Name)
	if err == nil {
		s.Value = value
		s.loaded = true
	}
	return err
}

func (s *state) Delete() error {
	err := s.client.DeleteParameter(s.Path + s.Name)
	s.Name = ""
	s.Value = ""
	s.loaded = false
	s.saved = false
	return err
}

// //////////////////////////////////////////////////////////////////////
// Service struct and methods
// //////////////////////////////////////////////////////////////////////
type Service struct {
	state state
}

func NewService(name string, value string, path string) (*Service, error) {
	state, err := NewState(name, value, path)
	if err != nil {
		return nil, err
	}
	return &Service{
		state: state,
	}, nil
}

func (s *Service) GetName() string {
	return s.state.GetName()
}

func (s *Service) GetValue() (string, error) {
	err := s.state.Load()
	if err != nil {
		return "", err
	}
	return s.state.GetValue(), nil

}

func (s *Service) SetValue(value string) error {
	s.state.SetValue(value)
	return s.state.Save()
}

func (s *Service) Delete() error {
	return s.state.Delete()
}

// ///////////////////////////////////////////////////////////////////////////
// Vault struct and methods
// ///////////////////////////////////////////////////////////////////////////
type Vault struct {
	services map[string]*Service
	ssmPath  string
}

// If no services provided, load all existing services from backend
func NewVault(services map[string]*Service, ssmPath string) (*Vault, error) {
	vault := &Vault{
		services: services,
		ssmPath:  ssmPath,
	}
	if services == nil {
		vault.services = make(map[string]*Service)
		err := vault.FindServices()
		if err != nil {
			return nil, err
		}
	}
	return vault, nil
}

func (v *Vault) Add(name string, value string) error {
	_, exists := v.services[name]
	if exists {
		return fmt.Errorf("service %s already exists", name)
	}
	service, err := NewService(name, value, v.ssmPath)
	if err != nil {
		return err
	}
	v.services[name] = service
	return nil
}

func (v *Vault) Get(name string) (string, error) {
	value, err := v.services[name].GetValue()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (v *Vault) Replace(name string, value string) error {
	service, exists := v.services[name]
	if !exists {
		return fmt.Errorf("service %s does not exist", name)
	}
	err := service.SetValue(value)
	return err
}

func (v *Vault) Delete(name string) error {
	service, exists := v.services[name]
	if !exists {
		return fmt.Errorf("service %s does not exist", name)
	}
	err := service.Delete()
	if err != nil {
		return err
	}
	delete(v.services, name)
	return nil
}

func (v *Vault) ListServices() []string {
	names := make([]string, 0, len(v.services))
	for name := range v.services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (v *Vault) FindServices() error {
	ssmClient := NewSsmClient()
	parameters, err := ssmClient.GetAllParametersByPath(v.ssmPath)
	if err != nil {
		return err
	}

	// Load services in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make(chan error, len(parameters))

	for name, value := range parameters {
		wg.Add(1)
		go func(name, value string) {
			defer wg.Done()
			service, err := NewService(name, value, v.ssmPath)
			if err != nil {
				errors <- err
				return
			}
			mu.Lock()
			v.services[name] = service
			mu.Unlock()
		}(name, value)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return <-errors
	}

	return nil
}
