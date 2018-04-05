package filter

import "unicode"

type Censor struct {
	filter *Filter
}

func NewCensor(filter *Filter) *Censor {
	return &Censor{filter}
}

func (censor *Censor) CleanSentence(inputSentence string) (sentence string) {
	var pass bool

	wordStart := 0
	inWord := false
	hasPeriod := false
	checkWord := []rune{}
	sentenceSlice := append([]rune(inputSentence), ' ')

	for i := 0; i < len(sentenceSlice); i++ {
		if unicode.IsSpace(sentenceSlice[i]) {
			if inWord {
				if hasPeriod {
					pass = censor.filter.TestRolling(checkWord)
				} else {
					pass = censor.filter.TestGeneral(checkWord)
				}
				if !pass {
					sentenceSlice = replaceSlice(sentenceSlice, '*', wordStart, i)
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
			checkWord = append(checkWord, unicode.ToLower(sentenceSlice[i]))
		}
	}
	return string(sentenceSlice[:len(sentenceSlice)-1])
}
