package filter

import "unicode"


type replacement struct {
	character rune
	modifier float64
}

func newReplacement(character rune, modifier float64) *replacement {
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

/*  Creates a new filter.  That large list is the 1000 most common words in English.  This
 *  This greatly reduces execution time by providing a large bank of whitelisted terms.  The
 *  most expensive operation in the filter is scoring a word.  Even though it takes time to
 *  iterate over this list, and this line is ridiculously long, it improves execution time
 *  by a factor of 2 in some instances.
 */
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
		known: map[string]bool {"the": true, "of": true, "to": true, "and": true, "a": true, "in": true, "is": true, "it": true, "you": true, "that": true, "he": true, "was": true, "for": true, "on": true, "are": true, "with": true, "as": true, "I": true, "his": true, "they": true, "be": true, "at": true, "one": true, "have": true, "this": true, "from": true, "or": true, "had": true, "by": true, "hot": true, "word": true, "but": true, "what": true, "some": true, "we": true, "can": true, "out": true, "other": true, "were": true, "all": true, "there": true, "when": true, "up": true, "use": true, "your": true, "how": true, "said": true, "an": true, "each": true, "she": true, "which": true, "do": true, "their": true, "time": true, "if": true, "will": true, "way": true, "about": true, "many": true, "then": true, "them": true, "write": true, "would": true, "like": true, "so": true, "these": true, "her": true, "long": true, "make": true, "thing": true, "see": true, "him": true, "two": true, "has": true, "look": true, "more": true, "day": true, "could": true, "go": true, "come": true, "did": true, "number": true, "sound": true, "no": true, "most": true, "people": true, "my": true, "over": true, "know": true, "water": true, "than": true, "call": true, "first": true, "who": true, "may": true, "down": true, "side": true, "been": true, "now": true, "find": true, "any": true, "new": true, "work": true, "part": true, "take": true, "get": true, "place": true, "made": true, "live": true, "where": true, "after": true, "back": true, "little": true, "only": true, "round": true, "man": true, "year": true, "came": true, "show": true, "every": true, "good": true, "me": true, "give": true, "our": true, "under": true, "name": true, "very": true, "through": true, "just": true, "form": true, "sentence": true, "great": true, "think": true, "say": true, "help": true, "low": true, "line": true, "differ": true, "turn": true, "cause": true, "much": true, "mean": true, "before": true, "move": true, "right": true, "boy": true, "old": true, "too": true, "same": true, "tell": true, "does": true, "set": true, "three": true, "want": true, "air": true, "well": true, "also": true, "play": true, "small": true, "end": true, "put": true, "home": true, "read": true, "hand": true, "port": true, "large": true, "spell": true, "add": true, "even": true, "land": true, "here": true, "must": true, "big": true, "high": true, "such": true, "follow": true, "act": true, "why": true, "ask": true, "men": true, "change": true, "went": true, "light": true, "kind": true, "off": true, "need": true, "house": true, "picture": true, "try": true, "us": true, "again": true, "animal": true, "point": true, "mother": true, "world": true, "near": true, "build": true, "self": true, "earth": true, "father": true, "head": true, "stand": true, "own": true, "page": true, "should": true, "country": true, "found": true, "answer": true, "school": true, "grow": true, "study": true, "still": true, "learn": true, "plant": true, "cover": true, "food": true, "sun": true, "four": true, "between": true, "state": true, "keep": true, "eye": true, "never": true, "last": true, "let": true, "thought": true, "city": true, "tree": true, "cross": true, "farm": true, "hard": true, "start": true, "might": true, "story": true, "saw": true, "far": true, "sea": true, "draw": true, "left": true, "late": true, "run": true, "don't": true, "while": true, "press": true, "close": true, "night": true, "real": true, "life": true, "few": true, "north": true, "open": true, "seem": true, "together": true, "next": true, "white": true, "children": true, "begin": true, "got": true, "walk": true, "example": true, "ease": true, "paper": true, "group": true, "always": true, "music": true, "those": true, "both": true, "mark": true, "often": true, "letter": true, "until": true, "mile": true, "river": true, "car": true, "feet": true, "care": true, "second": true, "book": true, "carry": true, "took": true, "science": true, "eat": true, "room": true, "friend": true, "began": true, "idea": true, "fish": true, "mountain": true, "stop": true, "once": true, "base": true, "hear": true, "horse": true, "cut": true, "sure": true, "watch": true, "color": true, "face": true, "wood": true, "main": true, "enough": true, "plain": true, "girl": true, "usual": true, "young": true, "ready": true, "above": true, "ever": true, "red": true, "list": true, "though": true, "feel": true, "talk": true, "bird": true, "soon": true, "body": true, "dog": true, "family": true, "direct": true, "pose": true, "leave": true, "song": true, "measure": true, "door": true, "product": true, "black": true, "short": true, "numeral": true, "class": true, "wind": true, "question": true, "happen": true, "complete": true, "ship": true, "area": true, "half": true, "rock": true, "order": true, "fire": true, "south": true, "problem": true, "piece": true, "told": true, "knew": true, "pass": true, "since": true, "top": true, "whole": true, "king": true, "space": true, "heard": true, "best": true, "hour": true, "better": true, "true": true, "during": true, "hundred": true, "five": true, "remember": true, "step": true, "early": true, "hold": true, "west": true, "ground": true, "interest": true, "reach": true, "fast": true, "verb": true, "sing": true, "listen": true, "six": true, "table": true, "travel": true, "less": true, "morning": true, "ten": true, "simple": true, "several": true, "vowel": true, "toward": true, "war": true, "lay": true, "against": true, "pattern": true, "slow": true, "center": true, "love": true, "person": true, "money": true, "serve": true, "appear": true, "road": true, "map": true, "rain": true, "rule": true, "govern": true, "pull": true, "cold": true, "notice": true, "voice": true, "unit": true, "power": true, "town": true, "fine": true, "certain": true, "fly": true, "fall": true, "lead": true, "cry": true, "dark": true, "machine": true, "note": true, "wait": true, "plan": true, "figure": true, "star": true, "box": true, "noun": true, "field": true, "rest": true, "correct": true, "able": true, "pound": true, "done": true, "beauty": true, "drive": true, "stood": true, "contain": true, "front": true, "teach": true, "week": true, "final": true, "gave": true, "green": true, "oh": true, "quick": true, "develop": true, "ocean": true, "warm": true, "free": true, "minute": true, "strong": true, "special": true, "mind": true, "behind": true, "clear": true, "tail": true, "produce": true, "fact": true, "street": true, "inch": true, "multiply": true, "nothing": true, "course": true, "stay": true, "wheel": true, "full": true, "force": true, "blue": true, "object": true, "decide": true, "surface": true, "deep": true, "moon": true, "island": true, "foot": true, "system": true, "busy": true, "test": true, "record": true, "boat": true, "common": true, "gold": true, "possible": true, "plane": true, "stead": true, "dry": true, "wonder": true, "laugh": true, "thousand": true, "ago": true, "ran": true, "check": true, "game": true, "shape": true, "equate": true, "miss": true, "brought": true, "heat": true, "snow": true, "tire": true, "bring": true, "yes": true, "distant": true, "fill": true, "east": true, "paint": true, "language": true, "among": true, "grand": true, "ball": true, "yet": true, "wave": true, "drop": true, "heart": true, "am": true, "present": true, "heavy": true, "dance": true, "engine": true, "position": true, "arm": true, "wide": true, "sail": true, "material": true, "size": true, "vary": true, "settle": true, "speak": true, "weight": true, "general": true, "ice": true, "matter": true, "circle": true, "pair": true, "include": true, "divide": true, "syllable": true, "felt": true, "perhaps": true, "pick": true, "sudden": true, "count": true, "square": true, "reason": true, "length": true, "represent": true, "art": true, "subject": true, "region": true, "energy": true, "hunt": true, "probable": true, "bed": true, "brother": true, "egg": true, "ride": true, "cell": true, "believe": true, "fraction": true, "forest": true, "sit": true, "race": true, "window": true, "store": true, "summer": true, "train": true, "sleep": true, "prove": true, "lone": true, "leg": true, "exercise": true, "wall": true, "catch": true, "mount": true, "wish": true, "sky": true, "board": true, "joy": true, "winter": true, "sat": true, "written": true, "wild": true, "instrument": true, "kept": true, "glass": true, "grass": true, "cow": true, "job": true, "edge": true, "sign": true, "visit": true, "past": true, "soft": true, "fun": true, "bright": true, "gas": true, "weather": true, "month": true, "million": true, "bear": true, "finish": true, "happy": true, "hope": true, "flower": true, "clothe": true, "strange": true, "gone": true, "jump": true, "baby": true, "eight": true, "village": true, "meet": true, "root": true, "buy": true, "raise": true, "solve": true, "metal": true, "whether": true, "push": true, "seven": true, "paragraph": true, "third": true, "shall": true, "held": true, "hair": true, "describe": true, "cook": true, "floor": true, "either": true, "result": true, "burn": true, "hill": true, "safe": true, "cat": true, "century": true, "consider": true, "type": true, "law": true, "bit": true, "coast": true, "copy": true, "phrase": true, "silent": true, "tall": true, "sand": true, "soil": true, "roll": true, "temperature": true, "finger": true, "industry": true, "value": true, "fight": true, "lie": true, "beat": true, "excite": true, "natural": true, "view": true, "sense": true, "ear": true, "else": true, "quite": true, "broke": true, "case": true, "middle": true, "kill": true, "son": true, "lake": true, "moment": true, "scale": true, "loud": true, "spring": true, "observe": true, "child": true, "straight": true, "consonant": true, "nation": true, "dictionary": true, "milk": true, "speed": true, "method": true, "organ": true, "pay": true, "age": true, "section": true, "dress": true, "cloud": true, "surprise": true, "quiet": true, "stone": true, "tiny": true, "climb": true, "cool": true, "design": true, "poor": true, "lot": true, "experiment": true, "bottom": true, "key": true, "iron": true, "single": true, "stick": true, "flat": true, "twenty": true, "skin": true, "smile": true, "crease": true, "hole": true, "trade": true, "melody": true, "trip": true, "office": true, "receive": true, "row": true, "mouth": true, "exact": true, "symbol": true, "die": true, "least": true, "trouble": true, "shout": true, "except": true, "wrote": true, "seed": true, "tone": true, "join": true, "suggest": true, "clean": true, "break": true, "lady": true, "yard": true, "rise": true, "bad": true, "blow": true, "oil": true, "blood": true, "touch": true, "grew": true, "cent": true, "mix": true, "team": true, "wire": true, "cost": true, "lost": true, "brown": true, "wear": true, "garden": true, "equal": true, "sent": true, "choose": true, "fell": true, "fit": true, "flow": true, "fair": true, "bank": true, "collect": true, "save": true, "control": true, "decimal": true, "gentle": true, "woman": true, "captain": true, "practice": true, "separate": true, "difficult": true, "doctor": true, "please": true, "protect": true, "noon": true, "whose": true, "locate": true, "ring": true, "character": true, "insect": true, "caught": true, "period": true, "indicate": true, "radio": true, "spoke": true, "atom": true, "human": true, "history": true, "effect": true, "electric": true, "expect": true, "crop": true, "modern": true, "element": true, "hit": true, "student": true, "corner": true, "party": true, "supply": true, "bone": true, "rail": true, "imagine": true, "provide": true, "agree": true, "thus": true, "capital": true, "won't": true, "chair": true, "danger": true, "fruit": true, "rich": true, "thick": true, "soldier": true, "process": true, "operate": true, "guess": true, "necessary": true, "sharp": true, "wing": true, "create": true, "neighbor": true, "wash": true, "bat": true, "rather": true, "crowd": true, "corn": true, "compare": true, "poem": true, "string": true, "bell": true, "depend": true, "meat": true, "rub": true, "tube": true, "famous": true, "dollar": true, "stream": true, "fear": true, "sight": true, "thin": true, "triangle": true, "planet": true, "hurry": true, "chief": true, "colony": true, "clock": true, "mine": true, "tie": true, "enter": true, "major": true, "fresh": true, "search": true, "send": true, "yellow": true, "gun": true, "allow": true, "print": true, "dead": true, "spot": true, "desert": true, "suit": true, "current": true, "lift": true, "rose": true, "continue": true, "block": true, "chart": true, "hat": true, "sell": true, "success": true, "company": true, "subtract": true, "event": true, "particular": true, "deal": true, "swim": true, "term": true, "opposite": true, "wife": true, "shoe": true, "shoulder": true, "spread": true, "arrange": true, "camp": true, "invent": true, "cotton": true, "born": true, "determine": true, "quart": true, "nine": true, "truck": true, "noise": true, "level": true, "chance": true, "gather": true, "shop": true, "stretch": true, "throw": true, "shine": true, "property": true, "column": true, "molecule": true, "select": true, "wrong": true, "gray": true, "repeat": true, "require": true, "broad": true, "prepare": true, "salt": true, "nose": true, "plural": true, "anger": true, "claim": true, "continent": true, "oxygen": true, "sugar": true, "death": true, "pretty": true, "skill": true, "women": true, "season": true, "solution": true, "magnet": true, "silver": true, "thank": true, "branch": true, "match": true, "suffix": true, "especially": true, "fig": true, "afraid": true, "huge": true, "sister": true, "steel": true, "discuss": true, "forward": true, "similar": true, "guide": true, "experience": true, "score": true, "apple": false, "bought": true, "led": true, "pitch": true, "coat": true, "mass": true, "card": true, "band": true, "rope": true, "slip": true, "win": true, "dream": true, "evening": true, "condition": true, "feed": true, "tool": true, "total": true, "basic": true, "smell": true, "valley": true, "nor": true, "double": true, "seat": true, "arrive": true, "master": true, "track": true, "parent": true, "shore": true, "division": true, "sheet": true, "substance": true, "favor": true, "connect": true, "post": true, "spend": true, "chord": true, "fat": true, "glad": true, "original": true, "share": true, "station": true, "dad": true, "bread": true, "charge": true, "proper": true, "bar": true, "offer": true, "segment": true, "slave": true, "duck": true, "instant": true, "market": true, "degree": true, "populate": true, "chick": true, "dear": true, "enemy": true, "reply": true, "drink": true, "occur": true, "support": true, "speech": true, "nature": true, "range": true, "steam": true, "motion": true, "path": true, "liquid": true, "log": true, "meant": true, "quotient": true, "teeth": true, "shell": true, "neck": true},
	}
}

/* Inserts a word into the tree of banned words */
func (filter *Filter) BanWord(word string) {
	filter.tree.insertWord(word)
}

/* Inserts a word to the list of whitelisted terms. */
func (filter *Filter) WhitelistWord(word string) {
	if !filter.known[word] {
		filter.known[word] = true
	}
}

/* Updates the tolerance.  A lower tolerance increases the likelihood of a word being banned. */
func (filter *Filter) UpdateTolerance(tol float64) {
	filter.tolerance = tol
}

/*  This takes the input sentence and returns the cleaned form.  It does this by iterating
 *  over the input string until it detects a word.  It saves the index where it first locates
 *  a non-space character and continues on until it encounters a space.  If it finds a period
 *  or comma at any point in the middle of the word, it sets the hasPeriod flag.  Once it goes
 *  to score the word, it uses this flag to determine which test function to use.  Depending
 *  on the return, it either replaces all characters in the word with *'s or simply moves on
 *  to the next.
 */
func (filter *Filter) CleanSentence(inputSentence string, userPCS int) (sentence string, newPCS int) {
	var modifier float64
	var newChar rune
	var pass bool

	newPCS = userPCS
	wordStart := 0
	inWord := false
	hasPeriod := false
	checkWord := []rune{}
	sentenceSlice := append([]rune(inputSentence), ' ')

	for i := 0; i < len(sentenceSlice); i++ {
		if unicode.IsSpace(sentenceSlice[i]) {
			if inWord {
				if !filter.known[string(checkWord)] {
					if hasPeriod {
						pass = filter.testPeriod(checkWord, modifier, userPCS)
					} else {
						pass = filter.testGeneral(checkWord, modifier, userPCS)
					}
					if !pass {
						sentenceSlice = replaceSlice(sentenceSlice, '*', wordStart, i)
						newPCS += 1
					}
				}
				hasPeriod = false
				inWord = false
				checkWord = []rune{}
			}
		} else {
			if !inWord {
				inWord = true
				modifier = 0
				wordStart = i
			}
			if sentenceSlice[i] == '.' || sentenceSlice[i] == ',' {
				hasPeriod = true
			}
			newChar, modifier = filter.replaceSymbol(sentenceSlice[i], modifier)
			checkWord = append(checkWord, newChar)
		}
	}
	return string(sentenceSlice[:len(sentenceSlice)-1]), newPCS
}

/*  This takes in the character and replaces it if needed.  Common special character substitutes
 *  are replaced (@ -> a).  If a change takes place, the modifier is updated and returned.
 */
func (filter *Filter) replaceSymbol(char rune, modifier float64) (rune, float64) {
	tup, isSymbol := filter.symbols[byte(char)]
	if isSymbol {
		modifier += tup.modifier
		char = tup.character
	}
	return char, modifier
}

/*  The general test scores the entire input word. */
func (filter *Filter) testGeneral(testSlice []rune, modifier float64, userPCS int) bool {
	pass := true
	if !filter.known[string(testSlice)]{
		score := filter.tree.getScore(string(testSlice), modifier)
		if score > filter.tolerance-(float64(userPCS)/100) {
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
func (filter *Filter) testPeriod(input []rune, modifier float64, userPCS int) bool {
	score := 0.0
	pass := true
	testSlice := []rune{}
	for _, c := range input {
		if unicode.IsPunct(c) {
			modifier += .75
		} else {
			testSlice = append(testSlice, c)
		}
	}
	if !filter.known[string(testSlice)] {
		if len(testSlice) > 4 {
			for i := 0; i < len(testSlice)-4; i++ {
				score = maxFloat(score, filter.tree.getScore(string(testSlice[:i+4]), modifier))
				if score > filter.tolerance-(float64(userPCS)/100) {
					break
				}
			}
		} else {
			score = filter.tree.getScore(string(testSlice), modifier)
		}
		if score > filter.tolerance-(float64(userPCS)/100) {
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
