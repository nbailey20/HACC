package helpers

import (
	"regexp"
	"testing"
)

func TestPasswordGenerate(t *testing.T) {
	pass := GeneratePassword("any", "any", 12, 20)
	pass2 := GeneratePassword("any", "any", 12, 20)
	if pass == "" || pass2 == "" {
		t.Fatalf("Error generating password")
	}
}

func TestPasswordGenerateNotEmpty(t *testing.T) {
	pass := GeneratePassword("any", "any", 12, 20)
	if pass == "" {
		t.Error("Expected non-empty password")
	}
}

func TestPasswordGenerateUnique(t *testing.T) {
	pass1 := GeneratePassword("any", "any", 12, 20)
	pass2 := GeneratePassword("any", "any", 12, 20)
	if pass1 == pass2 {
		t.Error("Expected different passwords on consecutive calls")
	}
}

func TestPasswordGenerateLength(t *testing.T) {
	pass := GeneratePassword("any", "any", 12, 20)
	if len(pass) == 0 {
		t.Error("Expected password with length > 0")
	}
}

func TestPasswordGenerateMinLength(t *testing.T) {
	minLen := 15
	maxLen := 25
	pass := GeneratePassword("any", "any", minLen, maxLen)
	if len(pass) < minLen {
		t.Errorf("Expected password length >= %d, got %d", minLen, len(pass))
	}
}

func TestPasswordGenerateMaxLength(t *testing.T) {
	minLen := 15
	maxLen := 25
	pass := GeneratePassword("any", "any", minLen, maxLen)
	if len(pass) > maxLen {
		t.Errorf("Expected password length <= %d, got %d", maxLen, len(pass))
	}
}

func TestPasswordGenerateDigitsRequired(t *testing.T) {
	pass := GeneratePassword("required", "any", 20, 30)
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount == 0 {
		t.Errorf("Expected at least 1 digit char with 'required', got 0 in password: %s", pass)
	}
}

func TestPasswordGenerateDigitsForbidden(t *testing.T) {
	pass := GeneratePassword("forbidden", "any", 20, 30)
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount > 0 {
		t.Errorf("Expected 0 digit chars with 'forbidden', got %d in password: %s", digitCount, pass)
	}
}

func TestPasswordGenerateSpecialsRequired(t *testing.T) {
	pass := GeneratePassword("any", "required", 20, 30)
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount == 0 {
		t.Errorf("Expected at least 1 special char with 'required', got 0 in password: %s", pass)
	}
}

func TestPasswordGenerateSpecialsForbidden(t *testing.T) {
	pass := GeneratePassword("any", "forbidden", 20, 30)
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount > 0 {
		t.Errorf("Expected 0 special chars with 'forbidden', got %d in password: %s", specialCount, pass)
	}
}

func TestPasswordGenerateMultipleParams(t *testing.T) {
	pass := GeneratePassword("required", "required", 20, 30)

	// Check length
	if len(pass) < 20 || len(pass) > 30 {
		t.Errorf("Expected password length between 20 and 30, got %d", len(pass))
	}

	// Check digit count
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount == 0 {
		t.Errorf("Expected at least 1 digit char with 'required', got 0")
	}

	// Check special char count
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount == 0 {
		t.Errorf("Expected at least 1 special char with 'required', got 0")
	}
}

func TestPasswordGenerateDigitsAny(t *testing.T) {
	pass := GeneratePassword("any", "any", 12, 20)
	if pass == "" {
		t.Error("Expected non-empty password with digits 'any'")
	}
}

func TestPasswordBasic(t *testing.T) {
	pass := GeneratePassword("forbiden", "forbidden", 12, 20)
	if pass == "" {
		t.Error("Expected non-empty password with no digits or specials")
	}
}
