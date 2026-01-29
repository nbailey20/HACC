package vault

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/nbailey20/hacc/hacc/config"
	"github.com/nbailey20/hacc/hacc/helpers"
)

func TestCredential(t *testing.T) {
	path := "/hackyclient/test"
	client, err := NewSsmClient("")
	if err != nil {
		t.Fatalf("Error creating SSM client: %v", err)
	}
	s, err := NewCredential("testName", "initialVal", path, "", "obfkey", client)
	if err != nil {
		t.Fatalf("Error creating state: %v", err)
	}

	name := s.GetName()
	if name != "testName" {
		t.Errorf("Expected name 'testName', got '%s'", name)
	}
	value, err := s.GetValue()
	if err != nil {
		t.Fatalf("Error retrieving state value: %s", err)
	}
	if value != "initialVal" {
		t.Errorf("Expected value 'initial', got '%s'", value)
	}
	s.SetValue("updated")
	value, err = s.GetValue()
	if err != nil {
		t.Fatalf("Error retrieving state value: %s", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	err = s.Save()
	if err != nil {
		t.Errorf("Error saving state: %v", err)
	}

	s2, err := NewCredential("testName", "", path, "", "obfkey", client)
	if err != nil {
		t.Fatalf("Error creating state: %v", err)
	}
	s2.Load()
	value, err = s2.GetValue()
	if err != nil {
		t.Fatalf("Error retrieving state value: %s", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	err = s2.Delete()
	if err != nil {
		t.Errorf("Error deleting state: %v", err)
	}
}

func TestService(t *testing.T) {
	path := "/hackytest/service"
	client, err := NewSsmClient("")
	if err != nil {
		t.Fatalf("Error creating SSM client: %v", err)
	}
	svc, err := NewService("testService", map[string]string{"testUser": "initialVal"}, path, "", "obfkey", client)
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}

	users := svc.GetUsers("")
	if len(users) != 1 || users[0] != "testUser" {
		t.Errorf("Expected users ['testUser'], got %v", users)
	}

	err = svc.Add("bob", "nowaypass")
	if err != nil {
		t.Errorf("Error adding user bob: %v", err)
	}
	numUsers := svc.numUsers
	if numUsers != 2 {
		t.Errorf("Expected 2 users, got %d", numUsers)
	}
	if !svc.HasUser("bob") {
		t.Errorf("Expected user 'bob' to be added")
	}
	value, err := svc.GetValue("testUser")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "initialVal" {
		t.Errorf("Expected value 'initialVal', got '%s'", value)
	}

	err = svc.SetValue("testUser", "updated")
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}
	value, err = svc.GetValue("testUser")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}
	time.Sleep(2 * time.Second) // wait for eventual consistency

	svc2, err := NewService("testService", map[string]string{}, path, "", "obfkey", client)
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}
	value, err = svc2.GetValue("testUser")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}
	err = svc2.Delete("testUser")
	if err != nil {
		t.Errorf("Error deleting service: %v", err)
	}
	err = svc.Delete("testUser")
	if err == nil {
		t.Errorf("Expected Error deleting already deleted service: %v", err)
	}
	numUsers = svc2.numUsers
	if numUsers != 1 {
		t.Errorf("Expected 1 user after deletion, got %d", numUsers)
	}

	err = svc.DeleteAll()
	if err != nil {
		t.Errorf("Error deleting all users: %v", err)
	}
	numUsers = svc.numUsers
	if numUsers != 0 {
		t.Errorf("Expected 0 users after DeleteAll, got %d", numUsers)
	}
}

func TestVaultGeneral(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test",
		ObfuscationKey: "obfkey",
	}
	vault, err := NewVault(nil, cfg)
	if err != nil {
		t.Fatalf("Error creating vault: %v", err)
	}
	err = vault.Add("service1", "user1", "value1")
	if err != nil {
		t.Fatalf("Error adding service1: %v", err)
	}
	err = vault.Add("service1", "user12", "value12")
	if err != nil {
		t.Fatalf("Error adding service1: %v", err)
	}
	err = vault.Add("service2", "user2", "value2")
	if err != nil {
		t.Fatalf("Error adding service2: %v", err)
	}
	if !vault.HasService("service1") {
		t.Errorf("Expected vault to have service1")
	}
	if !vault.HasService("service2") {
		t.Errorf("Expected vault to have service2")
	}
	users, err := vault.GetUsersForService("service1")
	if err != nil {
		t.Fatalf("Error getting users for service1: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users for service1")
	}
	err = vault.FindServices()
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}
	value, err := vault.Get("service1", "user1")
	if err != nil {
		t.Fatalf("Error getting service1: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1' for service1, got '%s'", value)
	}
	value, err = vault.Get("service2", "user2")
	if err != nil {
		t.Fatalf("Error getting service2: %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected value 'value2' for service2, got '%s'", value)
	}
	err = vault.Replace("service1", "user12", "newvalue12")
	if err != nil {
		t.Fatalf("Error replacing service1: %v", err)
	}
	value, err = vault.Get("service1", "user12")
	if err != nil {
		t.Fatalf("Error getting service1: %v", err)
	}
	if value != "newvalue12" {
		t.Errorf("Expected value 'newvalue12' for service1/user12, got '%s'", value)
	}
	value2 := vault.ListServices("serv")
	if len(value2) != 2 {
		t.Errorf("Expected 2 services, got %d", len(value2))
	}

	err = vault.Delete("service1", "user1")
	if err != nil {
		t.Fatalf("Error deleting service1: %v", err)
	}
	time.Sleep(2 * time.Second) // wait for eventual consistency

	vault2, err2 := NewVault(nil, cfg)
	if err2 != nil {
		t.Fatalf("Error creating vault2: %v", err2)
	}
	value3 := vault2.ListServices("servi")
	if len(value3) != 2 {
		t.Errorf("Expected 2 services, got %d", len(value3))
	}

	err = vault.Delete("service2", "user2")
	if err != nil {
		t.Fatalf("Error deleting service2: %v", err)
	}
	time.Sleep(3 * time.Second) // wait for eventual consistency
	vault3, err3 := NewVault(nil, cfg)
	if err3 != nil {
		t.Fatalf("Error creating vault3: %v", err3)
	}
	value4 := vault3.ListServices("servi")
	if len(value4) != 1 {
		t.Errorf("Expected 1 service, got %d", len(value4))
	}
	err = vault3.Delete("service1", "user12")
	if err != nil {
		t.Fatalf("Error deleting service1 user12: %v", err)
	}
	time.Sleep(3 * time.Second) // wait for eventual consistency
	value5 := vault3.ListServices("")
	if len(value5) != 0 {
		t.Errorf("Expected 0 services, got %d", len(value5))
	}
}

func TestVaultMultiAdd(t *testing.T) {
	cfg := config.AWSConfig{
		ParamPath:      "/hackyclient/test",
		ObfuscationKey: "obfkey",
	}
	vault, err := NewVault(nil, cfg)
	if err != nil {
		t.Fatalf("Error creating vault: %v", err)
	}
	var creds []helpers.FileCred
	for i := range 20 {
		x := rand.Intn(i+1) + 1
		creds = append(
			creds,
			helpers.FileCred{
				Username: fmt.Sprintf("bob%d", i),
				Password: fmt.Sprintf("bobpass%d", i),
				Service:  fmt.Sprintf("bobserv%d", x),
			},
		)
	}

	t.Cleanup(func() {
		time.Sleep(2 * time.Second) // eventual consistency

		for _, c := range creds {
			if err := vault.Delete(c.Service, c.Username); err != nil {
				t.Errorf("cleanup failed for %s/%s: %v", c.Service, c.Username, err)
			}
		}
	})

	results := vault.AddMulti(creds)
	for _, r := range results {
		if !r.Success {
			t.Errorf("Expected success = True adding multiple services, got false")
		}
		if r.Err != nil {
			t.Errorf("Expected error = nil adding multiple services, got %v", r.Err)
		}
	}
	for _, c := range creds {
		if !vault.HasService(c.Service) {
			t.Errorf("Expected service %s to be in Vault", c.Service)
		}
		p, err := vault.Get(c.Service, c.Username)
		if err != nil {
			t.Errorf("Expected no error retrieving pass for %s/%s: %v", c.Service, c.Username, err)
		}
		if p != c.Password {
			t.Errorf("Expected password %s for %s/%s", c.Password, c.Service, c.Username)
		}
	}
}
