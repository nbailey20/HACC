package helpers

import (
	"regexp"
	"testing"
)

func TestPasswordGenerate(t *testing.T) {
	pass := GeneratePassword(0, 0, 0, 30)
	pass2 := GeneratePassword(0, 0, 0, 30)
	if pass == "" || pass2 == "" {
		t.Fatalf("Error generating password")
	}
}

func TestPasswordGenerateNotEmpty(t *testing.T) {
	pass := GeneratePassword(0, 0, 12, 20)
	if pass == "" {
		t.Error("Expected non-empty password")
	}
}

func TestPasswordGenerateUnique(t *testing.T) {
	pass1 := GeneratePassword(0, 0, 12, 20)
	pass2 := GeneratePassword(0, 0, 12, 20)
	if pass1 == pass2 {
		t.Error("Expected different passwords on consecutive calls")
	}
}

func TestPasswordGenerateLength(t *testing.T) {
	pass := GeneratePassword(0, 0, 12, 20)
	if len(pass) == 0 {
		t.Error("Expected password with length > 0")
	}
}

func TestPasswordGenerateMinLength(t *testing.T) {
	minLen := 15
	maxLen := 25
	pass := GeneratePassword(0, 0, minLen, maxLen)
	if len(pass) < minLen {
		t.Errorf("Expected password length >= %d, got %d", minLen, len(pass))
	}
}

func TestPasswordGenerateMaxLength(t *testing.T) {
	minLen := 15
	maxLen := 25
	pass := GeneratePassword(0, 0, minLen, maxLen)
	if len(pass) > maxLen {
		t.Errorf("Expected password length <= %d, got %d", maxLen, len(pass))
	}
}

func TestPasswordGenerateMinDigitChars(t *testing.T) {
	minDigits := 2
	pass := GeneratePassword(minDigits, 0, 20, 30)
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount < minDigits {
		t.Errorf("Expected at least %d digit chars, got %d in password: %s", minDigits, digitCount, pass)
	}
}

func TestPasswordGenerateMinSpecialChars(t *testing.T) {
	minSpecial := 2
	pass := GeneratePassword(0, minSpecial, 20, 30)
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount < minSpecial {
		t.Errorf("Expected at least %d special chars, got %d in password: %s", minSpecial, specialCount, pass)
	}
}

func TestPasswordGenerateMultipleParams(t *testing.T) {
	minDigits := 2
	minSpecial := 1
	minLen := 20
	maxLen := 30
	pass := GeneratePassword(minDigits, minSpecial, minLen, maxLen)

	// Check length
	if len(pass) < minLen || len(pass) > maxLen {
		t.Errorf("Expected password length between %d and %d, got %d", minLen, maxLen, len(pass))
	}

	// Check digit count
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount < minDigits {
		t.Errorf("Expected at least %d digit chars, got %d", minDigits, digitCount)
	}

	// Check special char count
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount < minSpecial {
		t.Errorf("Expected at least %d special chars, got %d", minSpecial, specialCount)
	}
}

func TestPasswordGenerateZeroMinimums(t *testing.T) {
	pass := GeneratePassword(0, 0, 10, 15)
	if pass == "" {
		t.Error("Expected non-empty password even with zero minimums")
	}
}

func TestPasswordGenerateHighMinimums(t *testing.T) {
	minDigits := 3
	minSpecial := 3
	minLen := 20
	maxLen := 30
	pass := GeneratePassword(minDigits, minSpecial, minLen, maxLen)

	// Check length
	if len(pass) < minLen || len(pass) > maxLen {
		t.Errorf("Expected password length between %d and %d, got %d", minLen, maxLen, len(pass))
	}

	// Check digit count
	digitCount := 0
	for _, char := range pass {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}
	if digitCount < minDigits {
		t.Errorf("Expected at least %d digit chars, got %d", minDigits, digitCount)
	}

	// Check special char count
	specialPattern := regexp.MustCompile(`[@!$]`)
	specialCount := len(specialPattern.FindAllString(pass, -1))
	if specialCount < minSpecial {
		t.Errorf("Expected at least %d special chars, got %d", minSpecial, specialCount)
	}
}
