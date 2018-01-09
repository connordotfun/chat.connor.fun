package filter

import (
	"unicode"
)

const baseConfidence = .75

type Filter struct {
	tree *MetricTree
	userPCS float64
	toleranceModifier float64
}

/*  Creates a new filter.  That large list is the 1000 most common words in English.  This
 *  This greatly reduces execution time by providing a large bank of whitelisted terms.  The
 *  most expensive operation in the filter is scoring a word.  Even though it takes time to
 *  iterate over this list, and this line is ridiculously long, it improves execution time
 *  by a factor of 2 in some instances.
 */
func NewFilter(bannedTree *MetricTree) *Filter {
	return &Filter{
		tree: bannedTree,
		userPCS: 0,
	}
}

/* Inserts a word into the tree of banned words */
func (filter *Filter) BanWord(word string) {
	filter.tree.insertWord(word)
}

/*  This takes the input sentence and returns the cleaned form.  It does this by iterating
 *  over the input string until it detects a word.  It saves the index where it first locates
 *  a non-space character and continues on until it encounters a space.  If it finds a period
 *  or comma at any point in the middle of the word, it sets the hasPeriod flag.  Once it goes
 *  to score the word, it uses this flag to determine which test function to use.  Depending
 *  on the return, it either replaces all characters in the word with *'s or simply moves on
 *  to the next.
 */
func (filter *Filter) CleanSentence(inputSentence string) (sentence string) {
	var newChar rune
	var pass bool

	wordStart := 0
	inWord := false
	hasPeriod := false
	checkWord := []rune{}
	sentenceSlice := append([]rune(inputSentence), ' ')

	for i := 0; i < len(sentenceSlice); i++ {
		if unicode.IsSpace(sentenceSlice[i]) {
			if inWord {
				if !filter.tree.wordSet.isKnown(checkWord) {
					if hasPeriod {
						pass = filter.testPeriod(checkWord)
					} else {
						pass = filter.testGeneral(checkWord)
					}
					if !pass {
						sentenceSlice = replaceSlice(sentenceSlice, '*', wordStart, i)
						filter.userPCS += 1
					}
				}
				hasPeriod = false
				inWord = false
				checkWord = []rune{}
			}
		} else {
			if !inWord {
				inWord = true
				wordStart = i
			}
			if sentenceSlice[i] == '.' || sentenceSlice[i] == ',' {
				hasPeriod = true
			}
			newChar = filter.tree.wordSet.replaceSymbol(sentenceSlice[i])
			checkWord = append(checkWord, unicode.ToLower(newChar))
		}
	}
	return string(sentenceSlice[:len(sentenceSlice)-1])
}

/*  This takes in the character and replaces it if needed.  Common special character substitutes
 *  are replaced (@ -> a).  If a change takes place, the modifier is updated and returned.
 */

/*  The general test scores the entire input word. */
func (filter *Filter) testGeneral(testSlice []rune) bool {
	pass := true
	if !filter.tree.wordSet.isKnown(testSlice) {
		score := filter.tree.getScore(string(testSlice), 0)
		if score > (baseConfidence + filter.toleranceModifier)-(filter.userPCS/100) {
			pass = false
		}
	}
	return pass
}

/*  The period test first removes all of the periods and commas.  Whenever a removal takes place,
 *  it assesses a large penalty to the modifier.  A rolling selection of four characters from the
 *  edited rune slice is then tested.  If any tested subset receives a high enough score, the entire
 *  unedited string is failed.
 */
func (filter *Filter) testPeriod(input []rune) bool {
	score := 0.0
	modifier := 0.0
	pass := true
	testSlice := []rune{}
	for _, c := range input {
		if unicode.IsPunct(c) {
			modifier += .7
		} else {
			testSlice = append(testSlice, c)
		}
	}
	if !filter.tree.wordSet.isKnown(testSlice) {
		if len(testSlice) > 4 {
			for i := 0; i < len(testSlice)-4; i++ {
				score = maxFloat(score, filter.tree.getScore(string(testSlice[:i+4]), modifier))
				if score > (baseConfidence + filter.toleranceModifier)-(filter.userPCS/100) {
					break
				}
			}
		} else {
			score = filter.tree.getScore(string(testSlice), modifier)
		}
		if score > (baseConfidence + filter.toleranceModifier)-(filter.userPCS/100) {
			pass = false
		}
	}
	return pass
}

/* Replaces a subset of a slice with the same character. */
func replaceSlice(in []rune, char rune, start int, end int) []rune {
	for i := start; i < end; i++ {
		in[i] = char
	}
	return in
}
