package filter

import "testing"

func BenchmarkRead(b *testing.B) {
	filter := NewFilter(.85, []string{"the","of","and","a","to","in","is","you","that","it","he","was","for","on","are","as","with","his","they","I","at","be","this","have","from","or","one","had","by","word","but","not","what","all","were","we","when","your","can","said","there","use","an","each","which","she","do","how","their","if","will","up","other","about","out","many","then","them","these","so","some","her","would","make","like","him","into","time","has","look","two","more","write","go","see","number","no","way","could","people","my","than","first","water","been","call","who","oil","its","now","find","long","down","day","did","get","come","made","may","part"})

	for i := 0; i < b.N; i++ {
		filter.CleanSentence("this is a very long message.  i dont exxpect them to be this long in the future but you never know, yknow?", 0)
	}
}

func TestAddWord(t *testing.T) {
	filter := NewFilter(.85, []string{""})
	filter.AddWord("test")
	filter.AddWord("next")
	filter.AddWord("last")
}

func TestCleanSentence(t *testing.T) {
	filter := NewFilter(.85, []string{"apple", "orange"})
	actualClean, actualPCS := filter.CleanSentence("the @pple is orang3", 0)
	var expectedClean = "the *** is ***"
	var expectedPCS = 2

	if actualClean != expectedClean || actualPCS != expectedPCS {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
		t.Fatalf("Expected %s but got %s", expectedPCS, actualPCS)
	}
}

func TestCleanSentencePeriodSpace(t *testing.T) {
	filter := NewFilter(.85, []string{"apple", "orange"})
	actualClean, _:= filter.CleanSentence("apple..orange", 0)
	var expectedClean = "***"

	if actualClean != expectedClean {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
	}
}

func TestCleanSentencePeriodInCenter(t *testing.T) {
	filter := NewFilter(.85, []string{"apple", "orange"})
	actualClean, _:= filter.CleanSentence("@pp..le", 0)
	var expectedClean = "***"

	if actualClean != expectedClean {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
	}
}

func TestCleanSentenceLowPCS(t *testing.T) {
	filter := NewFilter(.85, []string{"apple", "orange"})
	actualClean, actualPCS := filter.CleanSentence("the @pple is orang3", 0)
	var expectedClean = "the *** is ***"
	var expectedPCS = 2

	if actualClean != expectedClean || actualPCS != expectedPCS {
		t.Fatalf("Expected %s but got %s", expectedClean, actualClean)
		t.Fatalf("Expected %s but got %s", expectedPCS, actualPCS)
	}
}

func TestReplaceAtIndex(t *testing.T) {
	actualResult := replaceAtIndex("tesf", 't', 3)
	var expectedResult = "test"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestGetScore(t *testing.T) {
	filter := NewFilter(.85, []string{"test"})

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
	filter := NewFilter(.85, []string{""})

	actualResult := filter.tree.distance("hello", "hello")
	var expectedResult = 0

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceDiffLength(t *testing.T) {
	filter := NewFilter(.85, []string{""})

	actualResult := filter.tree.distance("a", "toz")
	var expectedResult = 3

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}

func TestDistanceEmpty(t *testing.T) {
	filter := NewFilter(.85, []string{""})

	actualResult := filter.tree.distance("test", "")
	var expectedResult = 4

	if actualResult != expectedResult {
		t.Fatalf("Expected %d but got %d", expectedResult, actualResult)
	}
}