package vault

import (
	"fmt"
	"strings"
)

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
