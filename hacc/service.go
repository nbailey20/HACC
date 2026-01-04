package main

import (
	"fmt"
	"sort"
	"strings"
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
func NewState(name string, value string, path string) (*state, error) {
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
		return &state, err
	}
	return &state, nil
}

func (s *state) GetName() string {
	return s.Name
}

// Lazy backend fetching, don't pull the value before its needed
func (s *state) GetValue() (string, error) {
	if !s.loaded {
		err := s.Load()
		if err != nil {
			return "", err
		}
	}
	return s.Value, nil
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
	name     string
	state    map[string]*state
	numUsers int
	path     string
}

func NewService(serviceName string, credentials map[string]string, path string) (*Service, error) {
	svcState := make(map[string]*state)
	// if credentials map is empty, load existing users from backend
	if len(credentials) == 0 {
		svc := &Service{
			name:     serviceName,
			state:    svcState,
			numUsers: 0,
			path:     path,
		}
		err := svc.FindUsers()
		if err != nil {
			return nil, err
		}
		return svc, nil
	}

	for name, value := range credentials {
		st, err := NewState(name, value, path+serviceName+"/")
		if err != nil {
			return nil, err
		}
		svcState[name] = st
	}
	return &Service{
		name:     serviceName,
		state:    svcState,
		numUsers: len(credentials),
		path:     path,
	}, nil
}

func (s *Service) Add(username string, value string) error {
	_, exists := s.state[username]
	if exists {
		return fmt.Errorf("user %s already exists in service %s", username, s.name)
	}
	st, err := NewState(username, value, s.path+s.name+"/")
	if err != nil {
		return err
	}
	s.state[username] = st
	s.numUsers++
	return nil
}

func (s *Service) GetUsers(prefix string) []string {
	users := make([]string, 0, len(s.state))
	for username := range s.state {
		if strings.HasPrefix(username, prefix) {
			users = append(users, username)
		}
	}
	return users
}

func (s *Service) GetValue(username string) (string, error) {
	return s.state[username].GetValue()
}

func (s *Service) SetValue(username string, value string) error {
	s.state[username].SetValue(value)
	return s.state[username].Save()
}

func (s *Service) Delete(username string) error {
	err := s.state[username].Delete()
	if err == nil {
		delete(s.state, username)
		s.numUsers--
	}
	return err
}

func (s *Service) DeleteAll() error {
	var firstErr error

	for username := range s.state {
		err := s.Delete(username)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *Service) HasUser(username string) bool {
	_, exists := s.state[username]
	return exists
}

func (s *Service) FindUsers() error {
	ssmClient := NewSsmClient()
	// returns parameters with encrypted values
	parameters, err := ssmClient.GetAllParametersByPath(s.path + s.name + "/")
	if err != nil {
		return err
	}
	for name := range parameters {
		userName := strings.TrimPrefix(name, s.path+s.name+"/")
		// create state with empty value to pull from backend
		st, err := NewState(userName, "", s.path+s.name+"/")
		if err != nil {
			return err
		}
		s.state[userName] = st
		s.numUsers++
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////////////
// Vault struct and methods
// ///////////////////////////////////////////////////////////////////////////
type Vault struct {
	services map[string]*Service
	path     string
}

// If no services provided, load all existing services from backend
func NewVault(services map[string]*Service, path string) (*Vault, error) {
	vault := &Vault{
		services: services,
		path:     path,
	}
	if services == nil {
		vault.services = make(map[string]*Service)
		// load existing services from backend in-place
		err := vault.FindServices()
		if err != nil {
			return nil, err
		}
	}
	return vault, nil
}

func (v *Vault) Add(serviceName string, username string, value string) error {
	_, exists := v.services[serviceName]
	if !exists {
		service, err := NewService(serviceName, map[string]string{username: value}, v.path)
		if err != nil {
			return err
		}
		v.services[serviceName] = service
		return nil
	}
	return v.services[serviceName].Add(username, value)
}

func (v *Vault) Get(serviceName string, username string) (string, error) {
	service, exists := v.services[serviceName]
	if !exists {
		return "", fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.GetValue(username)
}

func (v *Vault) Replace(serviceName string, username string, value string) error {
	service, exists := v.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.SetValue(username, value)
}

func (v *Vault) Delete(serviceName string, username string) error {
	service, exists := v.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s does not exist", serviceName)
	}
	err := service.Delete(username)
	if err != nil {
		return err
	}
	if service.numUsers == 0 {
		delete(v.services, serviceName)
	}
	return nil
}

func (v *Vault) ListServices(prefix string) []string {
	// return sorted list of service names with prefix filtering
	// provide empty string for no prefix filtering
	names := make([]string, 0, len(v.services))
	for name := range v.services {
		if strings.HasPrefix(name, prefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func (v *Vault) FindServices() error {
	ssmClient := NewSsmClient()
	// returns parameters with encrypted values
	// parameters are of the form /path/serviceName/username
	parameters, err := ssmClient.GetAllParametersByPath(v.path)
	if err != nil {
		return err
	}

	// extract unique service names from parameters
	serviceNames := make([]string, 0)
	for name := range parameters {
		trimmed := strings.TrimPrefix(name, v.path)
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) == 2 {
			serviceNames = append(serviceNames, parts[0])
		}
	}

	// Load services in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make(chan error, len(serviceNames))

	for _, serviceName := range serviceNames {
		wg.Add(1)
		go func(serviceName string) {
			defer wg.Done()
			service, err := NewService(
				serviceName,
				map[string]string{}, // leave empty to pull from backend, don't overwrite with encrypted value
				v.path,
			)
			if err != nil {
				errors <- err
				return
			}
			mu.Lock()
			v.services[serviceName] = service
			mu.Unlock()
		}(serviceName)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return <-errors
	}

	return nil
}

func (v *Vault) HasService(serviceName string) bool {
	_, exists := v.services[serviceName]
	return exists
}
