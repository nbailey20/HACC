package vault

import (
	"fmt"
	"strings"
)

const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

// Helper functions to (de)obfuscate a SSM Parameter path name using a key
// and ensure the output only contains characters in the SSM Parameter name set
// preserves / characters so GetParametersByPath still works

func obfuscatePath(path, name, key string) (string, error) {
	obf_name, err := obfuscate(name, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", path, obf_name), nil
}

func obfuscate(plaintext, key string) (string, error) {
	if len(key) == 0 {
		return "", fmt.Errorf("key must not be empty")
	}

	parts := strings.Split(plaintext, "/")
	outParts := make([]string, 0, len(parts))
	for _, part := range parts {
		// Preserve empty segments (leading or double slashes)
		if part == "" {
			outParts = append(outParts, "")
			continue
		}
		cipher, err := obfuscateSegment(part, key)
		if err != nil {
			return "", err
		}
		outParts = append(outParts, cipher)
	}
	return strings.Join(outParts, "/"), nil
}

func obfuscateSegment(plaintext, key string) (string, error) {
	// map of allowed chars to their indices
	allowedIndex := make(map[byte]uint8, 64)
	for i := range len(allowed) {
		allowedIndex[allowed[i]] = uint8(i)
	}

	// Convert plaintext to a bit stream
	var bits []uint8
	for i := 0; i < len(plaintext); i++ {
		b := plaintext[i]
		for j := 7; j >= 0; j-- {
			// each b to the right j times and mask with 1
			// gets the jth bit of b
			bits = append(bits, (b>>j)&1)
		}
	}
	// Pad with zeros to make length divisible by 6
	for len(bits)%6 != 0 {
		bits = append(bits, 0)
	}

	out := make([]byte, 0, len(bits)/6)
	keyIndex := 0
	for i := 0; i < len(bits); i += 6 {
		// Pack 6 bits into a value
		var v uint8
		for j := range 6 {
			// shift sequence left and add next bit
			v = (v << 1) | bits[i+j]
		}

		// Get key value
		kc := key[keyIndex%len(key)]
		keyVal, ok := allowedIndex[kc]
		if !ok {
			return "", fmt.Errorf("key contains invalid character: %q", kc)
		}
		keyIndex++

		// XOR and map to allowed charset
		out = append(out, allowed[v^keyVal])
	}

	return string(out), nil
}

func deobfuscate(ciphertext, key string) (string, error) {
	if len(key) == 0 {
		return "", fmt.Errorf("key must not be empty")
	}

	parts := strings.Split(ciphertext, "/")
	outParts := make([]string, 0, len(parts))
	for _, part := range parts {
		// Preserve empty segments (leading or double slashes)
		if part == "" {
			outParts = append(outParts, "")
			continue
		}
		cipher, err := deobfuscateSegment(part, key)
		if err != nil {
			return "", err
		}
		outParts = append(outParts, cipher)
	}
	return strings.Join(outParts, "/"), nil
}

func deobfuscateSegment(ciphertext, key string) (string, error) {

	// Map allowed chars to 6-bit values
	allowedIndex := make(map[rune]uint8, len(allowed))
	for i, r := range allowed {
		allowedIndex[r] = uint8(i)
	}

	// Step 1: convert chars → 6-bit values → bit stream
	var bits []uint8
	for i, r := range ciphertext {
		val, ok := allowedIndex[r]
		if !ok {
			return "", fmt.Errorf("invalid character in ciphertext: %q", r)
		}

		keyChar := key[i%len(key)]
		keyVal := allowedIndex[rune(keyChar)]

		v := val ^ keyVal

		// expand 6 bits
		for j := 5; j >= 0; j-- {
			bits = append(bits, (v>>j)&1)
		}
	}

	// Step 2: reassemble bits into bytes
	var out []byte
	for i := 0; i+7 < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b = (b << 1) | bits[i+j]
		}
		out = append(out, b)
	}
	return string(out), nil
}
