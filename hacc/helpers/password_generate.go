package helpers

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type passSubResult struct {
	subsMade int
	password string
}

type password struct {
	numDigit   int
	numSpecial int
	value      string
}

func GeneratePassword(
	digits string, // any, required, forbidden
	specials string,
	minLen int,
	maxLen int,
) string {
	maxCharSwaps := 5 // Max number of character substitutions in password
	if maxLen == 0 {  // Default max value from input parsing
		maxLen = 30
	}
	numWordsInPass := 4   // XKCD-style passwords, works best with maxLen >= 11
	maxAttempts := 100000 // Try this many times to generate passwords which meet the min/max requirements
	digitCharMap := map[rune]rune{
		'b': '6',
		'e': '3',
		'i': '1',
		'o': '0',
		't': '7',
	}
	specialCharMap := map[rune]rune{
		'a': '@',
		'i': '!',
		's': '$',
	}
	result := "Did not find password meeting all requirements, try again."
	for range maxAttempts {
		generated := randomWordsWithSwaps(
			numWordsInPass,
			maxCharSwaps,
			digitCharMap,
			specialCharMap,
		)
		// if the generated password doesn't meet our requirements,
		// keep trying until we get one that does
		if len(generated.value) >= minLen &&
			len(generated.value) <= maxLen &&
			validSwaps(generated.numDigit, digits) &&
			validSwaps(generated.numSpecial, specials) {
			result = generated.value
			break
		}
	}
	return result
}

func validSwaps(numSwaps int, swapRequirement string) bool {
	switch swapRequirement {
	case "any":
		return true
	case "required":
		return numSwaps > 0
	case "forbidden":
		return numSwaps == 0
	}
	return false
}

func randomWordsWithSwaps(
	numWords int,
	maxSwaps int,
	digitCharMap map[rune]rune,
	specialCharMap map[rune]rune,
) password {
	var words []string
	var result password
	wordListSize := len(wordList)

	for range numWords {
		// get random element of wordList
		maxNum := big.NewInt(int64(wordListSize))
		lineNum, _ := rand.Int(rand.Reader, maxNum)
		words = append(words, wordList[lineNum.Int64()])
	}

	// CamelCase words
	for idx, word := range words {
		words[idx] = strings.ToUpper(word[:1]) + word[1:]
	}

	rawPassword := strings.Join(words, "")

	// Swap random # of letters with special/digit chars
	maxNum := big.NewInt(int64(maxSwaps))
	numCharSwaps64, _ := rand.Int(rand.Reader, maxNum)
	if numCharSwaps64.Int64() == 0 {
		numCharSwaps64 = big.NewInt(numCharSwaps64.Int64() + 1)
	}
	numDigitSwaps64, _ := rand.Int(rand.Reader, numCharSwaps64)
	numCharSwaps := int(numCharSwaps64.Int64())
	numDigitSwaps := int(numDigitSwaps64.Int64())
	numSpecialSwaps := numCharSwaps - numDigitSwaps

	numSubPassword := subChars(&rawPassword, digitCharMap, numDigitSwaps)
	result.numDigit = numSubPassword.subsMade
	// If not enough digit subs available, try to make extra special subs
	if numSubPassword.subsMade < numDigitSwaps {
		numSpecialSwaps += numDigitSwaps - numSubPassword.subsMade
	}

	// If not enough special subs available, oh well
	specialSubPassword := subChars(
		&numSubPassword.password,
		specialCharMap,
		numSpecialSwaps,
	)
	result.numSpecial = specialSubPassword.subsMade
	result.value = specialSubPassword.password
	return result
}

func subChars(
	pass *string,
	charMap map[rune]rune,
	numSubs int,
) passSubResult {
	var passRunes []rune
	for _, char := range *pass {
		passRunes = append(passRunes, char)
	}
	numPassChars := len(passRunes)
	subsMade := 0
	maxAttempts := 1000
	numAttempts := 0

	for subsMade < numSubs && numAttempts < maxAttempts {
		// choose random char in password
		max := big.NewInt(int64(numPassChars))
		randomInt64, _ := rand.Int(rand.Reader, max)
		randomInt := int(randomInt64.Int64())
		randomChar := passRunes[randomInt]

		for char, _ := range charMap {

			if randomChar == char {
				passRunes[randomInt] = charMap[randomChar]
				subsMade += 1
				break
			}
		}
		numAttempts += 1
	}

	newPass := string(passRunes)
	return passSubResult{
		password: newPass,
		subsMade: subsMade,
	}
}
