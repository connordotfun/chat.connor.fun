package filter

import (
	"testing"
)

func BenchmarkCleanLong(b *testing.B) { // Executes in .9 ms
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	for i := 0; i < b.N; i++ {
		censor.CleanSentence("this is a very long message.  i dont exxpect them to be much longer than this is")
	}
}

func TestCleanSentence(t *testing.T) {
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	actualClean := censor.CleanSentence("the sample is orang3")
	var expectedClean = "the ****** is ******"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestOrange(t *testing.T) {
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	actualClean := censor.CleanSentence("apple")
	var expectedClean = "*****"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestCleanSentencePeriodSpace(t *testing.T) {
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	actualClean := censor.CleanSentence("apple..orange")
	var expectedClean = "*************"

	if actualClean != expectedClean {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
	}
}

func TestCleanSentencePeriodInCenter(t *testing.T) {
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	actualClean := censor.CleanSentence("or.a..ng3")
	var expectedClean = "*********"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestRollingCombo(t *testing.T) {
	censor := Censor{NewFilter(NewTree("../assets/testBannedList.txt"))}
	actualClean := censor.CleanSentence("te.s")
	var expectedClean = "****"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestGetScore(t *testing.T) {
	filter := NewFilter(NewTree("../assets/testBannedList.txt"))
	actualResult := filter.tree.distance("test", "")
	var expectedResult = 4

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestMinIntSlice(t *testing.T) {
	actualResult := minIntSlice([]int{4, 2, 6, 12, 6, 9})
	var expectedResult = 2

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestDistanceHello(t *testing.T) {
	filter := NewFilter(NewTree("../assets/testBannedList.txt"))
	actualResult := filter.tree.distance("hello", "hello")
	var expectedResult = 0

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceDiffLength(t *testing.T) {
	filter := NewFilter(NewTree("../assets/testBannedList.txt"))
	actualResult := filter.tree.distance("a", "toz")
	var expectedResult = 3

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceEmpty(t *testing.T) {
	filter := NewFilter(NewTree("../assets/testBannedList.txt"))
	actualResult := filter.tree.distance("test", "")
	var expectedResult = 4

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestNewWordSet(t *testing.T) {
	set := newWordSet()
	set.WhitelistWord("lakdjsf")
}