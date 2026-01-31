package vault

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type service struct {
	name           string
	credentials    map[string]*credential
	numUsers       int
	path           string
	kmsId          string
	obfuscationKey string
	client         *ssmClient
	mu             sync.Mutex
}

func NewService(
	serviceName string,
	credentials map[string]string,
	path string,
	kmsId string,
	obfuscationKey string,
	client *ssmClient,
) (*service, error) {
	serviceCreds := make(map[string]*credential)
	// if credentials map is empty, load existing users from backend
	if len(credentials) == 0 {
		svc := &service{
			name:           serviceName,
			credentials:    serviceCreds,
			numUsers:       0,
			path:           path,
			kmsId:          kmsId,
			obfuscationKey: obfuscationKey,
			client:         client,
		}
		err := svc.FindUsers()
		if err != nil {
			return nil, err
		}
		return svc, nil
	}

	for name, value := range credentials {
		obfPath, err := obfuscatePath(path, serviceName, obfuscationKey)
		if err != nil {
			return nil, err
		}
		st, err := NewCredential(
			name,
			value,
			obfPath,
			kmsId,
			obfuscationKey,
			client,
		)
		if err != nil {
			return nil, err
		}
		serviceCreds[name] = st
	}
	return &service{
		name:           serviceName,
		credentials:    serviceCreds,
		numUsers:       len(credentials),
		path:           path,
		kmsId:          kmsId,
		obfuscationKey: obfuscationKey,
		client:         client,
	}, nil
}

func (s *service) Add(username string, value string) error {
	// ensure Add is a thread-safe method
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.credentials[username]
	if exists {
		return fmt.Errorf("user %s already exists in service %s", username, s.name)
	}
	obfPath, err := obfuscatePath(s.path, s.name, s.obfuscationKey)
	if err != nil {
		return err
	}
	st, err := NewCredential(
		username,
		value,
		obfPath,
		s.kmsId,
		s.obfuscationKey,
		s.client,
	)
	if err != nil {
		return err
	}

	s.credentials[username] = st
	s.numUsers++
	return nil
}

func (s *service) GetUsers(prefix string) []string {
	users := make([]string, 0, s.numUsers)
	for username := range s.credentials {
		if strings.HasPrefix(username, prefix) {
			users = append(users, username)
		}
	}
	sort.Strings(users)
	return users
}

func (s *service) GetValue(username string) (string, error) {
	if _, exists := s.credentials[username]; !exists {
		return "", fmt.Errorf("user %s does not exist in service %s", username, s.name)
	}
	return s.credentials[username].GetValue()
}

func (s *service) SetValue(username string, value string) error {
	s.credentials[username].SetValue(value)
	return s.credentials[username].Save()
}

func (s *service) NumUsers() int {
	return s.numUsers
}

func (s *service) Delete(username string) error {
	// Step 1: grab credential safely
	s.mu.Lock()
	cred, exists := s.credentials[username]
	s.mu.Unlock()

	if !exists {
		return nil // or ParameterNotFound equivalent
	}

	// Step 2: do network call WITHOUT holding lock
	err := cred.Delete()

	// Step 3: update in-memory state safely
	if err == nil || strings.Contains(err.Error(), "ParameterNotFound") {
		s.mu.Lock()
		// re-check in case another goroutine deleted it
		if _, stillExists := s.credentials[username]; stillExists {
			delete(s.credentials, username)
			s.numUsers--
		}
		s.mu.Unlock()
	}

	return err
}

func (s *service) DeleteAll() error {
	var firstErr error

	for username := range s.credentials {
		err := s.Delete(username)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (s *service) HasUser(username string) bool {
	_, exists := s.credentials[username]
	return exists
}

func (s *service) FindUsers() error {
	obfPath, err := obfuscatePath(s.path, s.name, s.obfuscationKey)
	if err != nil {
		return err
	}

	// returns parameters with encrypted values
	parameters, err := s.client.GetAllParametersByPath(obfPath + "/")
	if err != nil {
		return err
	}
	for fullName := range parameters {
		obfUserName := strings.TrimPrefix(
			fullName,
			obfPath+"/",
		)

		userName, err := deobfuscate(obfUserName, s.obfuscationKey)
		if err != nil {
			return err
		}

		st, err := NewCredential(
			userName,
			"", // create credential with empty value to pull from backend
			obfPath,
			s.kmsId,
			s.obfuscationKey,
			s.client,
		)
		if err != nil {
			return err
		}
		s.credentials[userName] = st
		s.numUsers++
	}
	return nil
}
