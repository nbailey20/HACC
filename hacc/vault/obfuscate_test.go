package vault

import "testing"

func TestObfuscate(t *testing.T) {
	plaintext := "/my/secret/parameter"
	key := "mysecretkey"
	obfuscated, err := obfuscate(plaintext, key)
	if err != nil {
		t.Fatalf("obfuscate failed: %v", err)
	}
	expected := "/my/qUhNE3rN/qEx2AhX2xcn-"
	if obfuscated != expected {
		t.Errorf("obfuscate returned %s, expected %s", obfuscated, expected)
	}
}

func TestDeobfuscate(t *testing.T) {
	obfuscated := "/my/qUhNE3rN/qEx2AhX2xcn-"
	key := "mysecretkey"
	deobfuscated, err := deobfuscate(obfuscated, key)
	if err != nil {
		t.Fatalf("deobfuscate failed: %v", err)
	}
	expected := "/my/secret/parameter"
	if deobfuscated != expected {
		t.Errorf("deobfuscate returned %s, expected %s", deobfuscated, expected)
	}
}
