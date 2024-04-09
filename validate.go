package agents

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//go:embed crawler-user-agents.json
var crawlersJson []byte

// Crawler contains information about one crawler.
type Crawler struct {
	// Regexp of User Agent of the crawler.
	Pattern string `json:"pattern"`

	// Discovery date.
	AdditionDate time.Time `json:"addition_date"`

	// Official url of the robot.
	URL string `json:"url"`

	// Examples of full User Agent strings.
	Instances []string `json:"instances"`
}

// Private type needed to convert addition_date from/to the format used in JSON.
type jsonCrawler struct {
	Pattern      string   `json:"pattern"`
	AdditionDate string   `json:"addition_date"`
	URL          string   `json:"url"`
	Instances    []string `json:"instances"`
}

const timeLayout = "2006/01/02"

func (c Crawler) MarshalJSON() ([]byte, error) {
	jc := jsonCrawler{
		Pattern:      c.Pattern,
		AdditionDate: c.AdditionDate.Format(timeLayout),
		URL:          c.URL,
		Instances:    c.Instances,
	}
	return json.Marshal(jc)
}

func (c *Crawler) UnmarshalJSON(b []byte) error {
	var jc jsonCrawler
	if err := json.Unmarshal(b, &jc); err != nil {
		return err
	}

	c.Pattern = jc.Pattern
	c.URL = jc.URL
	c.Instances = jc.Instances

	if c.Pattern == "" {
		return fmt.Errorf("empty pattern in record %s", string(b))
	}

	if jc.AdditionDate != "" {
		tim, err := time.ParseInLocation(timeLayout, jc.AdditionDate, time.UTC)
		if err != nil {
			return err
		}
		c.AdditionDate = tim
	}

	return nil
}

// The list of crawlers, built from contents of crawler-user-agents.json.
var Crawlers = func() []Crawler {
	var crawlers []Crawler
	if err := json.Unmarshal(crawlersJson, &crawlers); err != nil {
		panic(err)
	}
	return crawlers
}()

// analyzePattern expands a regular expression to the list of matching texts
// for plain search. The list is complete, i.e. iff a text matches the input
// pattern, then it contains at least one of the returned texts. If such a list
// can't be built, then the resulting list contains one element (main literal),
// it also returns built regexp object to run in this case. The main literal is
// a text that is contained in any matching text and is used to optimize search
// (pre-filter with this main literal before running a regexp). In the case such
// a main literal can't be found or the regexp is invalid, an error is returned.
func analyzePattern(pattern string) ([]string, *regexp.Regexp, error) {
	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return nil, nil, fmt.Errorf("re %q does not compile: %w", pattern, err)
	}
	re = re.Simplify()

	// Try to convert it to the list of literals.
	const maxLiterals = 100
	literals, ok := literalizeRegexp(re, maxLiterals)
	if ok {
		return literals, nil, nil
	}

	// Fallback to using a regexp, but we need some string serving as
	// an indicator of its possible presence.
	mainLiteral := findLongestCommonLiteral(re)
	const minLiteralLen = 3
	if len(mainLiteral) < minLiteralLen {
		return nil, nil, fmt.Errorf("re %q does not contain sufficiently long literal to serve an indicator. The longest literal is %q", pattern, mainLiteral)
	}

	return []string{mainLiteral}, regexp.MustCompile(pattern), nil
}

// literalizeRegexp expands a regexp to the list of matching sub-strings.
// Iff a text matches the regexp, it contains at least one of the returned
// texts. Argument maxLiterals regulates the maximum number of patterns to
// return. In case of an overflow or if it is impossible to build such a list
// from the regexp, false is returned.
func literalizeRegexp(re *syntax.Regexp, maxLiterals int) (literals []string, ok bool) {
	switch re.Op {
	case syntax.OpNoMatch:
		return nil, true

	case syntax.OpEmptyMatch:
		return []string{""}, true

	case syntax.OpLiteral:
		return unwrapCase(re, []string{string(re.Rune)}, maxLiterals)

	case syntax.OpCharClass:
		count := 0
		for i := 0; i < len(re.Rune); i += 2 {
			first := re.Rune[i]
			last := re.Rune[i+1]
			count += int(last - first + 1)
		}

		if count > maxLiterals {
			return nil, false
		}

		patterns := make([]string, 0, count)
		for i := 0; i < len(re.Rune); i += 2 {
			first := re.Rune[i]
			last := re.Rune[i+1]
			for r := first; r <= last; r++ {
				patterns = append(patterns, string([]rune{r}))
			}
		}

		return unwrapCase(re, patterns, maxLiterals)

	case syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		// Not supported.
		return nil, false

	case syntax.OpBeginLine, syntax.OpBeginText:
		return []string{"^"}, true

	case syntax.OpEndLine, syntax.OpEndText:
		return []string{"$"}, true

	case syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		// Not supported.
		return nil, false

	case syntax.OpCapture:
		subList, ok := literalizeRegexp(re.Sub[0], maxLiterals)
		if !ok {
			return nil, false
		}

		return unwrapCase(re, subList, maxLiterals)

	case syntax.OpStar, syntax.OpPlus:
		// Not supported.
		return nil, false

	case syntax.OpQuest:
		if re.Flags&syntax.FoldCase != 0 {
			return nil, false
		}

		subList, ok := literalizeRegexp(re.Sub[0], maxLiterals)
		if !ok {
			return nil, false
		}
		subList = append(subList, "")

		return subList, true

	case syntax.OpRepeat:
		// Not supported.
		return nil, false

	case syntax.OpConcat:
		if re.Flags&syntax.FoldCase != 0 {
			return nil, false
		}

		matrix := make([][]string, len(re.Sub))
		for i, sub := range re.Sub {
			subList, ok := literalizeRegexp(sub, maxLiterals)
			if !ok {
				return nil, false
			}
			matrix[i] = subList
		}

		return combinations(matrix, maxLiterals)

	case syntax.OpAlternate:
		results := []string{}
		for _, sub := range re.Sub {
			subList, ok := literalizeRegexp(sub, maxLiterals)
			if !ok {
				return nil, false
			}
			results = append(results, subList...)
		}

		if len(results) > maxLiterals {
			return nil, false
		}

		return unwrapCase(re, results, maxLiterals)

	default:
		// Not supported.
		return nil, false
	}
}

// combinations produces all combination of elements of matrix.
// Each sub-slice of matrix contributes one part of a resulting string.
// If the number of combinations is larger than maxLiterals, the function
// returns false.
func combinations(matrix [][]string, maxLiterals int) ([]string, bool) {
	if len(matrix) == 1 {
		if len(matrix[0]) > maxLiterals {
			return nil, false
		}

		return matrix[0], true
	}

	prefixes := matrix[0]
	suffixes, ok := combinations(matrix[1:], maxLiterals)
	if !ok {
		return nil, false
	}

	size := len(prefixes) * len(suffixes)
	if size > maxLiterals {
		return nil, false
	}

	results := make([]string, 0, size)
	for _, prefix := range prefixes {
		for _, suffix := range suffixes {
			results = append(results, prefix+suffix)
		}
	}

	return results, true
}

// unwrapCase takes the regexp and the list of patterns expanded from it and
// further expands it for a case-insensitive regexp, if needed. Argument
// maxLiterals regulates the maximum number of patterns to return. In case of an
// overflow, false is returned.
func unwrapCase(re *syntax.Regexp, patterns []string, maxLiterals int) ([]string, bool) {
	if re.Flags&syntax.FoldCase == 0 {
		return patterns, true
	}

	results := []string{}
	for _, pattern := range patterns {
		matrix := make([][]string, len(pattern))
		for i, r := range pattern {
			upper := unicode.ToUpper(r)
			lower := unicode.ToLower(r)
			matrix[i] = []string{
				string([]rune{upper}),
				string([]rune{lower}),
			}
		}

		patterns, ok := combinations(matrix, maxLiterals)
		if !ok {
			return nil, false
		}

		results = append(results, patterns...)
		if len(results) > maxLiterals {
			return nil, false
		}
	}

	return results, true
}

// findLongestCommonLiteral finds the longest common literal in the regexp. It's
// such a string which is contained in any text matching the regexp. If such a
// literal can't be found, it returns an empty string.
func findLongestCommonLiteral(re *syntax.Regexp) string {
	if re.Flags&syntax.FoldCase != 0 {
		return ""
	}

	switch re.Op {
	case syntax.OpNoMatch, syntax.OpEmptyMatch:
		return ""

	case syntax.OpLiteral:
		return string(re.Rune)

	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		return ""

	case syntax.OpBeginLine, syntax.OpBeginText:
		return "^"

	case syntax.OpEndLine, syntax.OpEndText:
		return "$"

	case syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return ""

	case syntax.OpCapture:
		return findLongestCommonLiteral(re.Sub[0])

	case syntax.OpStar:
		return ""

	case syntax.OpPlus:
		return findLongestCommonLiteral(re.Sub[0])

	case syntax.OpQuest:
		return ""

	case syntax.OpRepeat:
		if re.Min >= 1 {
			return findLongestCommonLiteral(re.Sub[0])
		}

		return ""

	case syntax.OpConcat:
		longest := ""
		for _, sub := range re.Sub {
			str := findLongestCommonLiteral(sub)
			if len(str) > len(longest) {
				longest = str
			}
		}

		return longest

	case syntax.OpAlternate:
		return ""

	default:
		return ""
	}
}

type regexpPattern struct {
	re    *regexp.Regexp
	index int
}

type matcher struct {
	replacer *strings.Replacer
	regexps  []regexpPattern
}

var uniqueToken = hex.EncodeToString((&maphash.Hash{}).Sum(nil))

const (
	uniqueTokenLen = 2 * 8
	numLen         = 5
	literalLabel   = '-'
	regexpLabel    = '*'
)

var m = func() matcher {
	if len(uniqueToken) != uniqueTokenLen {
		panic("len(uniqueToken) != uniqueTokenLen")
	}

	regexps := []regexpPattern{}
	oldnew := make([]string, 0, len(Crawlers)*2)

	// Put re-based patterns to the end to prevent AdsBot-Google from
	// shadowing AdsBot-Google-Mobile.
	var oldnew2 []string

	for i, crawler := range Crawlers {
		literals, re, err := analyzePattern(crawler.Pattern)
		if err != nil {
			panic(err)
		}

		label := literalLabel
		num := i
		if re != nil {
			label = regexpLabel
			num = len(regexps)
			regexps = append(regexps, regexpPattern{
				re:    re,
				index: i,
			})
		}

		replaceWith := fmt.Sprintf(" %s%c%0*d ", uniqueToken, label, numLen, num)

		for _, literal := range literals {
			if re != nil {
				oldnew2 = append(oldnew2, literal, replaceWith)
			} else {
				oldnew = append(oldnew, literal, replaceWith)
			}
		}
	}
	oldnew = append(oldnew, oldnew2...)

	// Allocate another array with regexps of exact size to save memory.
	regexps2 := make([]regexpPattern, len(regexps))
	copy(regexps2, regexps)

	r := strings.NewReplacer(oldnew...)
	r.Replace("") // To cause internal build process.

	return matcher{
		replacer: r,
		regexps:  regexps2,
	}
}()

// Returns if User Agent string matches any of crawler patterns.
func IsCrawler(userAgent string) bool {
	// This code is mostly copy-paste of MatchingCrawlers,
	// but with early exit logic, so it works a but faster.

	text := "^" + userAgent + "$"
	replaced := m.replacer.Replace(text)
	if replaced == text {
		return false
	}

	for {
		uniquePos := strings.Index(replaced, uniqueToken)
		if uniquePos == -1 {
			break
		}

		start := uniquePos + uniqueTokenLen + 1
		if start+numLen >= len(replaced) {
			panic("corrupt replaced: " + replaced)
		}

		label := replaced[start-1]
		switch label {
		case literalLabel:
			return true
		case regexpLabel:
			// Rare case. Run regexp to confirm the match.
			indexStr := replaced[start : start+numLen]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				panic("corrupt replaced: " + replaced)
			}
			rp := m.regexps[index]
			if rp.re.MatchString(userAgent) {
				return true
			}
		default:
			panic("corrupt replaced: " + replaced)
		}

		replaced = replaced[start+numLen:]
	}

	return false
}

// Finds all crawlers matching the User Agent and returns the list of their indices in Crawlers.
func MatchingCrawlers(userAgent string) []int {
	text := "^" + userAgent + "$"
	replaced := m.replacer.Replace(text)
	if replaced == text {
		return []int{}
	}

	indices := []int{}
	for {
		uniquePos := strings.Index(replaced, uniqueToken)
		if uniquePos == -1 {
			break
		}

		start := uniquePos + uniqueTokenLen + 1
		if start+numLen >= len(replaced) {
			panic("corrupt replaced: " + replaced)
		}
		indexStr := replaced[start : start+numLen]
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			panic("corrupt replaced: " + replaced)
		}

		label := replaced[start-1]
		switch label {
		case literalLabel:
			indices = append(indices, index)
		case regexpLabel:
			// Rare case. Run regexp to confirm the match.
			rp := m.regexps[index]
			if rp.re.MatchString(userAgent) {
				indices = append(indices, rp.index)
			}
		default:
			panic("corrupt replaced: " + replaced)
		}

		replaced = replaced[start+numLen:]
	}

	return indices
}
