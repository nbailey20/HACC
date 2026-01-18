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

func GeneratePassword() string {
	maxCharSwaps := 5   // Max number of character substitutions in password
	numWordsInPass := 4 // XKCD-style passwords
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

	return randomWordsWithSwaps(
		numWordsInPass,
		maxCharSwaps,
		digitCharMap,
		specialCharMap,
	)
	// Perform random character substitions in given password
}

func randomWordsWithSwaps(
	numWords int,
	maxSwaps int,
	digitCharMap map[rune]rune,
	specialCharMap map[rune]rune,
) string {
	var words []string
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

	password := strings.Join(words, "")

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

	subResult := subChars(&password, digitCharMap, numDigitSwaps)
	// If not enough digit subs available, try to make extra special subs
	if subResult.subsMade < numDigitSwaps {
		numSpecialSwaps += numDigitSwaps - subResult.subsMade
	}

	// If not enough special subs available, oh well
	return subChars(
		&subResult.password,
		specialCharMap,
		numSpecialSwaps,
	).password
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
