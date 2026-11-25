package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
)

type FileWordList struct {
	words  map[int]string
	digits int
}

func (w *FileWordList) Digits() int {
	return w.digits
}

func (w *FileWordList) WordAt(i int) string {
	word, ok := w.words[i]
	if !ok {
		return ""
	}
	return word
}

func (w *FileWordList) NumWords() int {
	return len(w.words)
}

func LoadWordListFromFile(path string) (*FileWordList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open wordlist file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	words := make(map[int]string)
	seenWords := make(map[string]struct{})
	digits := 0
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format in %s at line %d: %q", path, lineNo, line)
		}

		keyStr := parts[0]
		word := parts[1]

		key, err := strconv.Atoi(keyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid dice key in %s at line %d: %q", path, lineNo, keyStr)
		}

		if digits == 0 {
			digits = len(keyStr)
		} else if len(keyStr) != digits {
			return nil, fmt.Errorf("inconsistent key length in %s at line %d: got %d digits, expected %d",
				path, lineNo, len(keyStr), digits)
		}

		if strings.Trim(keyStr, "123456") != "" {
			return nil, fmt.Errorf("invalid dice key in %s at line %d: %q (digits must be 1-6)", path, lineNo, keyStr)
		}

		if _, exists := words[key]; exists {
			return nil, fmt.Errorf("duplicate key in %s at line %d: %d", path, lineNo, key)
		}

		if _, exists := seenWords[word]; exists {
			return nil, fmt.Errorf("duplicate word in %s at line %d: %q", path, lineNo, word)
		}
		seenWords[word] = struct{}{}

		words[key] = word
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan wordlist file: %w", err)
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("wordlist is empty")
	}

	// Every possible dice roll must map to a word. Otherwise WordAt returns ""
	// for missing keys and passphrases silently lose entropy.
	expected := 1
	for i := 0; i < digits; i++ {
		expected *= 6
	}
	if len(words) != expected {
		return nil, fmt.Errorf("incomplete wordlist %s: got %d words, expected %d for %d dice",
			path, len(words), expected, digits)
	}

	return &FileWordList{
		digits: digits,
		words:  words,
	}, nil
}

// Compile-time interface checks
var (
	_ diceware.WordList = (*FileWordList)(nil)
)
