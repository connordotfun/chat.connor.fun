package filter

import (
	"os"
	"bufio"
)

type node struct {
	word string
	children map[int]*node
}

func newNode(inputWord string) *node {
	return &node {
		word: inputWord,
		children: make(map[int]*node),
	}
}

type MetricTree struct {
	root *node
	tolerance int
	wordSet *wordSet
}

func NewTree(path string) *MetricTree {
	tree := &MetricTree{
		root: nil,
		tolerance: 3,
		wordSet: newWordSet(),
	}

	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineStr := scanner.Text()
		tree.insertWord(lineStr)
	}

	return tree
}

func (tree *MetricTree) insertWord(word string) {
	if tree.root == nil {
		tree.root = newNode(word)
	} else {
		tree.insertWordInternal(newNode(word))
	}
}

/* The "search" case if you will */
func (tree *MetricTree) insertWordInternal(node *node) {
	var parent = tree.root

	var dist = tree.distance(parent.word, node.word)

	for child, hasChild := parent.children[dist]; hasChild; {
		parent = child
		dist = tree.distance(parent.word, node.word)
		child, hasChild = parent.children[dist]
	}
	parent.children[dist] = node
}

/*  This calculates the Levenshtein distance between two words.  Think string alignment. */
func (tree *MetricTree) distance(s1, s2 string) int {
	s1Len := len(s1)
	s2Len := len(s2)

	if s1Len == 0 {
		return s2Len
	}
	if s2Len == 0 {
		return s1Len
	}
	dist := make([][]int, s1Len+1)
	for x := 0; x < s1Len+1; x++ {
		dist[x] = make([]int, s2Len+1)
	}

	for i := 1; i < s1Len+1; i++ {
		dist[i][0] = i
		for j:= 1; j < s2Len+1; j++ {
			var cost int
			if s1[i-1] == s2[j-1] {
				cost = 0
			} else {
				cost = 1
			}
			if i == 1 {
				dist[0][j] = j
			}
			dist[i][j] = min(min(dist[i-1][j]+1,
				dist[i][j-1]+1),
				dist[i-1][j-1]+cost)
			if i == s1Len && j == s2Len {
			}
			if i>1 && j>1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				dist[i][j] = min(dist[i][j], dist[i-2][j-2]+cost)
			}
		}
	}
	return dist[s1Len][s2Len]
}

/*  This scores the inputted word on how similar it is to any word in the database tree.
 *  This outer function calls the recursive search which generates a list of the shortest
 *  distances between the input and any close words in the tree.  This slice is then returned
 *  to this function.  The smallest value is used to create the score. The score the percent
 *  of the input that matches the closest word in the tree plus the input's modifier.
 */
func (tree *MetricTree) getScore(word string, modifier float64) (score float64){
	distList := []int{}
	resultStrings := tree.getScoreRec(word, tree.root)
	for i := 0; i < len(resultStrings); i++ {
		distList = append(distList, tree.distance(resultStrings[i], word))
	}
	if len(distList) > 0 {
		closest := minIntSlice(distList)
		score = float64((float64(len(word)) - float64(closest) + modifier)/float64(len(word)))
	} else {
		score = 0.0
	}
	return
}

/* The "search" case if you will */
func (tree *MetricTree) getScoreRec(word string, root *node) []string {
	dist := tree.distance(word, root.word)
	similarWords := []string{}
	if dist < tree.tolerance {
		similarWords = append(similarWords, root.word)
	}
	minVal := max(1, dist-tree.tolerance)
	maxVal := dist+tree.tolerance

	for i := minVal; i <= maxVal; i++ {
		if child, hasChild := root.children[i]; hasChild {
			similarWords = append(similarWords, tree.getScoreRec(word, child)...)
		}
	}
	return similarWords
}

/* Finds the minimum int in a slice of ints. */
func minIntSlice(v []int) (m int) {
	m = v[0]
	for _, e := range v {
		if e < m {
			m = e
		}
	}
	return
}

/* This functions are simple inequalites for various datatypes. */
func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func maxFloat(a, b float64) float64{
	if a > b {
		return a
	} else {
		return b
	}
}