package filter

import (
	"testing"
)

func BenchmarkCleanLong(b *testing.B) { // Executes in .9 ms
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	for i := 0; i < b.N; i++ {
		filter.CleanSentence("this is a very long message.  i dont exxpect them to be much longer than this is")
	}
}

func TestCleanSentence(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	//filter.BanWord("orange")
	filter.BanWord("sample")
	actualClean := filter.CleanSentence("the sample is orang3")
	var expectedClean = "the ****** is ******"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestOrange(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	filter.BanWord("apple")
	actualClean := filter.CleanSentence("apple")
	var expectedClean = "*****"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestCleanSentencePeriodSpace(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	filter.BanWord("orange")
	filter.BanWord("orange")
	actualClean := filter.CleanSentence("apple..orange")
	var expectedClean = "*************"

	if actualClean != expectedClean {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
	}
}

func TestCleanSentencePeriodInCenter(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	actualClean := filter.CleanSentence("or.a..ng3")
	var expectedClean = "*********"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestPeriodCombo(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	actualClean := filter.CleanSentence("te.s")
	var expectedClean = "****"

	if actualClean != expectedClean {
		t.Fatalf("Expected '%s' but got '%s'", expectedClean, actualClean)
	}
}

func TestGetScore(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
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
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	actualResult := filter.tree.distance("hello", "hello")
	var expectedResult = 0

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceDiffLength(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
	actualResult := filter.tree.distance("a", "toz")
	var expectedResult = 3

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceEmpty(t *testing.T) {
	filter := NewFilter(NewTree("../assets/bannedList.txt"))
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