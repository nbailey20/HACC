package helpers

import "testing"

func TestPasswordGenerate(t *testing.T) {
	pass := GeneratePassword()
	pass2 := GeneratePassword()
	if pass == "" || pass2 == "" {
		t.Fatalf("Error generating password")
	}
}

func TestPasswordGenerateNotEmpty(t *testing.T) {
	pass := GeneratePassword()
	if pass == "" {
		t.Error("Expected non-empty password")
	}
}

func TestPasswordGenerateUnique(t *testing.T) {
	pass1 := GeneratePassword()
	pass2 := GeneratePassword()
	if pass1 == pass2 {
		t.Error("Expected different passwords on consecutive calls")
	}
}

func TestPasswordGenerateLength(t *testing.T) {
	pass := GeneratePassword()
	if len(pass) == 0 {
		t.Error("Expected password with length > 0")
	}
}
