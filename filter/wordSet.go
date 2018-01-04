package filter

import (
	"bufio"
	"os"
)

type wordSet struct {
	whitelist map[string]bool
	symbols map[rune]rune
}

func newWordSet() *wordSet{
	set := wordSet{
		whitelist: make(map[string]bool),
		symbols: map[rune]rune{
			'@': 'a',
			'!': 'i',
			'3': 'e',
			'$': 's',
			'&': 'a',
		},
	}

	file, _ := os.Open("whitelist.txt")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr := scanner.Text()
		set.whitelist[string(lineStr)] = true
	}
	return &set
}

/* Checks if the input is whitelisted. */
func (set *wordSet) isKnown(test []rune) bool {
	return set.whitelist[string(test)]
}

/* Takes a symbol and replaces it if needed. */
func (set *wordSet) replaceSymbol(char rune) rune {
	replace, isSymbol := set.symbols[char]
	if isSymbol {
		char = replace
	}
	return char
}

/* Inserts a word to the list of whitelisted terms. */
func (set *wordSet) WhitelistWord(word string) {
	if !set.whitelist[word] {
		set.whitelist[word] = true
	}
}