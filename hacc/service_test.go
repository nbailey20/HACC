package main

import (
	"testing"
	"time"
)

func TestState(t *testing.T) {
	s, err := NewState("/hackyclient/test/state", "initial")
	if err != nil {
		t.Fatalf("Error creating state: %v", err)
	}

	name := s.GetName()
	if name != "/hackyclient/test/state" {
		t.Errorf("Expected name '/hackyclient/test/state', got '%s'", name)
	}
	value := s.GetValue()
	if value != "initial" {
		t.Errorf("Expected value 'initial', got '%s'", value)
	}
	s.SetValue("updated")
	value = s.GetValue()
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	err = s.Save()
	if err != nil {
		t.Errorf("Error saving state: %v", err)
	}

	s2, err := NewState("/hackyclient/test/state", "")
	if err != nil {
		t.Fatalf("Error creating state: %v", err)
	}
	s2.Load()
	value = s2.GetValue()
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	err = s2.Delete()
	if err != nil {
		t.Errorf("Error deleting state: %v", err)
	}
}

func TestService(t *testing.T) {
	svc, err := NewService("/hackyclient/test/service", "initial")
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}
	name := svc.GetName()
	if name != "/hackyclient/test/service" {
		t.Errorf("Expected name '/hackyclient/test/service', got '%s'", name)
	}

	value, err := svc.GetValue()
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "initial" {
		t.Errorf("Expected value 'initial', got '%s'", value)
	}

	err = svc.SetValue("updated")
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}
	value, err = svc.GetValue()
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	svc2, err := NewService("/hackyclient/test/service", "")
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}
	value, err = svc2.GetValue()
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if value != "updated" {
		t.Errorf("Expected value 'updated', got '%s'", value)
	}

	err = svc2.Delete()
	if err != nil {
		t.Errorf("Error deleting service: %v", err)
	}
	err = svc.Delete()
	if err != nil {
		t.Errorf("Error deleting already deleted service: %v", err)
	}
}

func TestVault(t *testing.T) {
	vault, err := NewVault(nil, "/hackyclient/test/")
	if err != nil {
		t.Fatalf("Error creating vault: %v", err)
	}
	err = vault.Add("/hackyclient/test/service1", "value1")
	if err != nil {
		t.Fatalf("Error adding service1: %v", err)
	}
	err = vault.Add("/hackyclient/test/service2", "value2")
	if err != nil {
		t.Fatalf("Error adding service2: %v", err)
	}
	err = vault.FindServices()
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}
	value, err := vault.Get("/hackyclient/test/service1")
	if err != nil {
		t.Fatalf("Error getting service1: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1' for service1, got '%s'", value)
	}
	value, err = vault.Get("/hackyclient/test/service2")
	if err != nil {
		t.Fatalf("Error getting service2: %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected value 'value2' for service2, got '%s'", value)
	}

	err = vault.Replace("/hackyclient/test/service1", "newvalue1")
	if err != nil {
		t.Fatalf("Error replacing service1: %v", err)
	}
	value, err = vault.Get("/hackyclient/test/service1")
	if err != nil {
		t.Fatalf("Error getting service1: %v", err)
	}
	if value != "newvalue1" {
		t.Errorf("Expected value 'newvalue1' for service1, got '%s'", value)
	}

	value2 := vault.ListServices()
	if len(value2) != 2 {
		t.Errorf("Expected 2 services, got %d", len(value2))
	}

	err = vault.Delete("/hackyclient/test/service1")
	if err != nil {
		t.Fatalf("Error deleting service1: %v", err)
	}
	time.Sleep(2 * time.Second) // wait for eventual consistency

	vault2, err2 := NewVault(nil, "/hackyclient/test/")
	if err2 != nil {
		t.Fatalf("Error creating vault2: %v", err2)
	}
	value3 := vault2.ListServices()
	if len(value3) != 1 {
		t.Errorf("Expected 1 services, got %d", len(value3))
	}

	err = vault.Delete("/hackyclient/test/service2")
	if err != nil {
		t.Fatalf("Error deleting service2: %v", err)
	}
}
