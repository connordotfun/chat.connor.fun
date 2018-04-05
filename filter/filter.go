package filter

const baseConfidence = .75

type Filter struct {
	tree *MetricTree
	toleranceModifier float64
	wordSet *wordSet
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
		wordSet: newWordSet(),
	}
}

/* Inserts a word into the tree of banned words */
func (filter *Filter) BanWord(word string) {
	filter.tree.insertWord(word)
}


/*  This takes in the character and replaces it if needed.  Common special character substitutes
 *  are replaced (@ -> a).  If a change takes place, the modifier is updated and returned.
 */

/*  The general test scores the entire input word. */
func (filter *Filter) TestGeneral(input []rune) bool {
	pass := true
	testSlice, _ := filter.replaceSymbols(input)
	if !filter.wordSet.isKnown(testSlice) {
		score := filter.tree.getScore(string(testSlice), 0)
		if score > (baseConfidence + filter.toleranceModifier) {
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
func (filter *Filter) TestRolling(input []rune) bool {
	score := 0.0
	pass := true
	testSlice, modifier := filter.replaceSymbols(input)
	if !filter.wordSet.isKnown(testSlice) {
		if len(testSlice) > 4 {
			for i := 0; i < len(testSlice)-4; i++ {
				score = maxFloat(score, filter.tree.getScore(string(testSlice[:i+4]), modifier))
				if score > (baseConfidence + filter.toleranceModifier) {
					break
				}
			}
		} else {
			score = filter.tree.getScore(string(testSlice), modifier)
		}
		if score > (baseConfidence + filter.toleranceModifier) {
			pass = false
		}
	}
	return pass
}

func (filter *Filter) replaceSymbols (input []rune) (newSlice []rune, modifier float64){
	newSlice = []rune{}
	modifier = 0
	for _, c := range input {
		if c == ',' || c == '.' {
			modifier += .5
		} else {
			newSlice = append(newSlice, filter.wordSet.replaceSymbol(c))
		}
	}
	return
}

/* Replaces a subset of a slice with the same character. */
func replaceSlice(in []rune, char rune, start int, end int) []rune {
	for i := start; i < end; i++ {
		in[i] = char
	}
	return in
}
