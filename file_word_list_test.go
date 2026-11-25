package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWordListFromFile(t *testing.T) {
	if _, err := LoadWordListFromFile("./wordlists/DE-wordlist-diceware.txt"); err != nil {
		t.Fatalf("bundled wordlist failed to load: %v", err)
	}

	cases := map[string]string{
		"incomplete":     "1\ta\n2\tb\n3\tc\n4\td\n5\te\n",
		"invalid digit":  "1\ta\n2\tb\n3\tc\n4\td\n5\te\n7\tf\n",
		"duplicate word": "1\ta\n2\tb\n3\tc\n4\td\n5\te\n6\ta\n",
	}
	for name, content := range cases {
		path := filepath.Join(t.TempDir(), "list.txt")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadWordListFromFile(path); err == nil {
			t.Fatalf("%s: expected an error, got nil", name)
		}
	}
}
