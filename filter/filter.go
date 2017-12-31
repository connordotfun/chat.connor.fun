package filter

import (
	"strings"
)

type replacement struct {
	character byte
	modifier float64
}

func newReplacement(character byte, modifier float64) *replacement {
	return &replacement{
		character: character,
		modifier: modifier,
	}
}

type Filter struct {
	tolerance float64
	tree *metricTree
	known map[string]bool
	symbols map[byte]*replacement
	banned []string
}

func NewFilter(tol float64, banned []string) *Filter {
	return &Filter{
		tolerance: tol,
		tree: newTree(banned),
		symbols: map[byte]*replacement{
			'@': newReplacement('a', -.2),
			'!': newReplacement('i', -.2),
			'3': newReplacement('e', -.2),
			'$': newReplacement('s', -.2),
		},
	}
}

func (filter *Filter) AddWord(word string) {
	filter.tree.insertWord(word)
}

func (filter *Filter) CleanSentence(inputSentence string, userPCS int) (sentence string, newPCS int) {
	words := strings.Fields(inputSentence)
	sentence = inputSentence
	newPCS = userPCS
	for _, word := range words {
		modifier := 0.0
		checkWord := strings.ToLower(word)
		if found := filter.known[word]; !found {
			for i := 0; i < len(word); i++ {
				tup, isSymbol := filter.symbols[word[i]]
				if isSymbol {
					modifier += tup.modifier
					checkWord = replaceAtIndex(checkWord, tup.character, i)
				}
			}
			checkWordDel := strings.Replace(checkWord, ".", "", -1)
			checkWordSpace := strings.Replace(checkWord, ".", " ", -1)

			var scoreDel, scoreSpace float64

			if checkWordDel == checkWordSpace {
				scoreDel = filter.tree.getScore(checkWordDel, modifier)
				scoreSpace = 0
			} else {
				scoreDel = filter.tree.getScore(checkWordDel, modifier)
				scoreSpace = filter.tree.getScore(checkWordSpace, modifier)
			}

			if maxFloat(scoreDel, scoreSpace) > filter.tolerance-(float64(userPCS)/100) {
				sentence = strings.Replace(sentence, word, "***", -1)
				newPCS += 1
			}
		}
	}
	return
}

func replaceAtIndex(in string, b byte, i int) string {
	out := []rune(in)
	out[i] = rune(b)
	return string(out)
}