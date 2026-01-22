package vault

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/nbailey20/hacc/hacc/config"
)

type Vault struct {
	Services map[string]*service
	path     string
	keyId    string
	client   *ssmClient
}

// If no services provided, load all existing services from backend
func NewVault(services map[string]*service, cfg config.AWSConfig) (*Vault, error) {
	client := NewSsmClient(cfg.Profile)
	vault := &Vault{
		Services: services,
		path:     cfg.ParamPath,
		keyId:    cfg.KmsId,
		client:   client,
	}
	if services == nil {
		vault.Services = make(map[string]*service)
		// load existing services from backend in-place
		err := vault.FindServices()
		if err != nil {
			return nil, err
		}
	}
	return vault, nil
}

func (v *Vault) Add(serviceName string, username string, value string) error {
	_, exists := v.Services[serviceName]
	if !exists {
		service, err := NewService(
			serviceName,
			map[string]string{username: value},
			v.path,
			v.keyId,
			v.client,
		)
		if err != nil {
			return err
		}
		v.Services[serviceName] = service
		return nil
	}
	return v.Services[serviceName].Add(username, value)
}

func (v *Vault) Get(serviceName string, username string) (string, error) {
	service, exists := v.Services[serviceName]
	if !exists {
		return "", fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.GetValue(username)
}

func (v *Vault) Replace(serviceName string, username string, value string) error {
	service, exists := v.Services[serviceName]
	if !exists {
		return fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.SetValue(username, value)
}

func (v *Vault) Delete(serviceName string, username string) error {
	service, exists := v.Services[serviceName]
	if !exists {
		return fmt.Errorf("service %s does not exist", serviceName)
	}
	err := service.Delete(username)
	if err != nil {
		return err
	}
	if service.numUsers == 0 {
		delete(v.Services, serviceName)
	}
	return nil
}

func (v *Vault) ListServices(prefix string) []string {
	// return sorted list of service names with prefix filtering
	// provide empty string for no prefix filtering
	names := make([]string, 0, len(v.Services))
	for name := range v.Services {
		if strings.HasPrefix(name, prefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func (v *Vault) FindServices() error {
	// returns parameters with encrypted values
	// parameters are of the form /path/serviceName/username
	parameters, err := v.client.GetAllParametersByPath(v.path)
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
				v.keyId,
				v.client,
			)
			if err != nil {
				errors <- err
				return
			}
			mu.Lock()
			v.Services[serviceName] = service
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
	_, exists := v.Services[serviceName]
	return exists
}
