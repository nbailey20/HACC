package vault

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nbailey20/hacc/config"
	"github.com/nbailey20/hacc/helpers"
)

var MAX_CONCURRENT_REQUESTS = 3

type Vault struct {
	Services       map[string]*service
	path           string
	keyId          string
	obfuscationKey string
	client         *ssmClient
	mu             sync.Mutex
}

type AddCredResult struct {
	Service  string
	Username string
	Success  bool
	Err      error
}

// If no services provided, load all existing services from backend
func NewVault(services map[string]*service, cfg config.AWSConfig) (*Vault, error) {
	client, err := NewSsmClient(cfg.Profile)
	vault := &Vault{
		Services:       services,
		path:           cfg.ParamPath,
		keyId:          cfg.KmsId,
		obfuscationKey: cfg.ObfuscationKey,
		client:         client,
	}
	if err != nil {
		return nil, err
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
	v.mu.Lock()
	_, exists := v.Services[serviceName]
	v.mu.Unlock()

	if !exists {
		service, err := NewService(
			serviceName,
			map[string]string{username: value},
			v.path,
			v.keyId,
			v.obfuscationKey,
			v.client,
		)
		if err != nil {
			return err
		}
		// check again to make sure service wasn't created since previous locked check
		v.mu.Lock()
		if _, exists := v.Services[serviceName]; !exists {
			v.Services[serviceName] = service
		} else {
			// service was created in the meantime, add credential to existing service
			err = v.Services[serviceName].Add(username, value)
		}
		v.mu.Unlock()
		return nil
	}
	return v.Services[serviceName].Add(username, value)
}

func (v *Vault) AddMulti(creds []helpers.FileCred) []AddCredResult {
	resultsChan := make(chan AddCredResult, len(creds))
	sem := make(chan struct{}, MAX_CONCURRENT_REQUESTS)
	var wg sync.WaitGroup

	for _, cred := range creds {
		wg.Add(1)
		go func(cred helpers.FileCred) {
			defer wg.Done()
			sem <- struct{}{} // acquire semaphore
			defer func() {
				if r := recover(); r != nil {
					resultsChan <- AddCredResult{
						Service:  cred.Service,
						Username: cred.Username,
						Success:  false,
						Err:      fmt.Errorf("panic: %v", r),
					}
				}
				<-sem // release semaphore
			}()
			// slight random sleep to reduce likelihood of throttling
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			err := v.Add(
				cred.Service,
				cred.Username,
				cred.Password,
			)
			result := AddCredResult{
				Service:  cred.Service,
				Username: cred.Username,
				Success:  true,
				Err:      err,
			}
			if err != nil {
				result.Success = false
			}
			resultsChan <- result
		}(cred)
	}

	// close results channel when all workers finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// read the results of the channel
	out := make([]AddCredResult, 0, len(creds))
	for r := range resultsChan {
		out = append(out, r)
	}
	return out
}

func (v *Vault) Get(serviceName string, username string) (string, error) {
	service, exists := v.Services[serviceName]
	if !exists {
		return "", fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.GetValue(username)
}

func (v *Vault) GetUsersForService(serviceName string) ([]string, error) {
	service, exists := v.Services[serviceName]
	if !exists {
		return []string{}, fmt.Errorf("service %s does not exist", serviceName)
	}
	return service.GetUsers(""), nil
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
	parameters, err := v.client.GetAllParametersByPath(v.path + "/")
	if err != nil {
		return err
	}

	// extract unique service names from parameters
	serviceNames := make([]string, 0)
	seen := make(map[string]bool)
	for fullName := range parameters {
		obfServiceUser := strings.TrimPrefix(fullName, v.path+"/")
		parts := strings.SplitN(obfServiceUser, "/", 2)
		service, err := deobfuscate(parts[0], v.obfuscationKey)
		if err != nil {
			return err
		}
		if len(parts) == 2 && !seen[service] {
			serviceNames = append(serviceNames, service)
			seen[service] = true
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
				v.obfuscationKey,
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
