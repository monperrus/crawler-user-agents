package agents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp/syntax"
	"sort"
	"strings"
	"testing"
)

// TestAnalyzePattern tests analyzePattern function on many cases, including
// edge cases.
func TestAnalyzePattern(t *testing.T) {
	cases := []struct {
		input            string
		wantError        string
		wantPatterns     []string
		wantRe           bool
		shouldMatchRe    []string
		shouldNotMatchRe []string
	}{
		{
			input:        "simple phrase",
			wantPatterns: []string{"simple phrase"},
		},
		{
			input:        "^begin anchor",
			wantPatterns: []string{"^begin anchor"},
		},
		{
			input:        "end anchor$",
			wantPatterns: []string{"end anchor$"},
		},
		{
			input:        "^both anchors$",
			wantPatterns: []string{"^both anchors$"},
		},
		{
			input:        "(alter|nation)",
			wantPatterns: []string{"alter", "nation"},
		},
		{
			input:            "too many [aA][lL][tT][eE][rR][nN][aA][tT][iI][oO][nN][sS]",
			wantPatterns:     []string{"too many "},
			wantRe:           true,
			shouldMatchRe:    []string{"too many ALTERNATIONs"},
			shouldNotMatchRe: []string{"too many combinations "},
		},
		{
			input: "(alter|nation) concatenation (alter|nation)",
			wantPatterns: []string{
				"alter concatenation alter",
				"alter concatenation nation",
				"nation concatenation alter",
				"nation concatenation nation",
			},
		},
		{
			input: "clas[sS] of [c]haract[eiu]rs",
			wantPatterns: []string{
				"clasS of characters",
				"clasS of charactirs",
				"clasS of characturs",
				"class of characters",
				"class of charactirs",
				"class of characturs",
			},
		},
		{
			input: "ranges [0-3]x[a-c]",
			wantPatterns: []string{
				"ranges 0xa", "ranges 0xb", "ranges 0xc",
				"ranges 1xa", "ranges 1xb", "ranges 1xc",
				"ranges 2xa", "ranges 2xb", "ranges 2xc",
				"ranges 3xa", "ranges 3xb", "ranges 3xc",
			},
		},
		{
			input:        "Quest?",
			wantPatterns: []string{"Ques", "Quest"},
		},
		{
			input:        "Q?ue(st)?",
			wantPatterns: []string{"Que", "Quest", "ue", "uest"},
		},
		{
			input:            "too many combinations [0-9][a-z]",
			wantPatterns:     []string{"too many combinations "},
			wantRe:           true,
			shouldMatchRe:    []string{"too many combinations 0a"},
			shouldNotMatchRe: []string{"too many combinations "},
		},
		{
			input:            "negation in char class [^x]",
			wantPatterns:     []string{"negation in char class "},
			wantRe:           true,
			shouldMatchRe:    []string{"negation in char class y"},
			shouldNotMatchRe: []string{"negation in char class x"},
		},
		{
			input:            "any char .",
			wantPatterns:     []string{"any char "},
			wantRe:           true,
			shouldMatchRe:    []string{"any char x"},
			shouldNotMatchRe: []string{"any char_x"},
		},
		{
			input:            `word \boundary`,
			wantPatterns:     []string{"oundary"},
			wantRe:           true,
			shouldMatchRe:    []string{"word oundary"},
			shouldNotMatchRe: []string{"word boundary"},
		},
		{
			input:            "asterisk*",
			wantPatterns:     []string{"asteris"},
			wantRe:           true,
			shouldMatchRe:    []string{"asteris", "asterisk", "asteriskk"},
			shouldNotMatchRe: []string{"asterik"},
		},
		{
			input:            "plus+",
			wantPatterns:     []string{"plu"},
			wantRe:           true,
			shouldMatchRe:    []string{"plus", "pluss"},
			shouldNotMatchRe: []string{"plu"},
		},
		{
			input:        "repeat{3,5}$",
			wantPatterns: []string{"repeattt$", "repeatttt$", "repeattttt$"},
		},
		{
			input:            "repeat{1,120}$",
			wantPatterns:     []string{"repea"},
			wantRe:           true,
			shouldMatchRe:    []string{"repeattt", "repeatttt", "repeattttt"},
			shouldNotMatchRe: []string{"repea5"},
		},
		{
			input:     "broken re[",
			wantError: "does not compile",
		},
		{
			input:     "n?o? ?l?o?n?g? ?l?i?t?e?r?a?l?",
			wantError: "does not contain sufficiently long literal",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			gotPatterns, re, err := analyzePattern(tc.input)
			if tc.wantError != "" {
				if err == nil {
					t.Fatalf("expected to get an error, got success")
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("the error returned must contain text %q, got %q", tc.wantError, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			sort.Strings(tc.wantPatterns)
			sort.Strings(gotPatterns)
			if !reflect.DeepEqual(tc.wantPatterns, gotPatterns) {
				t.Fatalf("returned list of patterns (%#v) does not match the expected value (%#v)", gotPatterns, tc.wantPatterns)
			}

			if !tc.wantRe {
				if re != nil {
					t.Fatalf("unexpectedly got a re")
				}

				return
			}

			if re == nil {
				t.Fatalf("expected to get a re, got nil")
			}
			for _, text := range tc.shouldMatchRe {
				if !re.MatchString(text) {
					t.Fatalf("test %q must match against the re, but it doesn't", text)
				}
			}
			for _, text := range tc.shouldNotMatchRe {
				if re.MatchString(text) {
					t.Fatalf("test %q must not match against the re, but it does", text)
				}
			}
		})
	}
}

// TestLiteralizeRegexp tests expansion of a regexp to a list of literals.
func TestLiteralizeRegexp(t *testing.T) {
	cases := []struct {
		input        string
		maxLiterals  int
		wantOutput   []string
		wantOverflow bool
	}{
		{
			input:       "simple phrase",
			maxLiterals: 100,
			wantOutput:  []string{"simple phrase"},
		},
		{
			input:       "cases [1-2x-z]",
			maxLiterals: 100,
			wantOutput:  []string{"cases 1", "cases 2", "cases x", "cases y", "cases z"},
		},
		{
			input:       "[Ii]gnore case",
			maxLiterals: 100,
			wantOutput:  []string{"Ignore case", "ignore case"},
		},
		{
			input:        "overflow [1-2x-z]",
			maxLiterals:  2,
			wantOverflow: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			re, err := syntax.Parse(tc.input, syntax.Perl)
			if err != nil {
				t.Fatalf("failed to parse regexp %q: %v", tc.input, err)
			}

			gotPatterns, ok := literalizeRegexp(re, tc.maxLiterals)
			if tc.wantOverflow {
				if ok {
					t.Fatalf("expected to get an overflow, got success")
				}

				return
			}

			if !ok {
				t.Fatalf("unexpected overflow")
			}

			sort.Strings(tc.wantOutput)
			sort.Strings(gotPatterns)
			if !reflect.DeepEqual(tc.wantOutput, gotPatterns) {
				t.Fatalf("returned list of patterns (%#v) does not match the expected value (%#v)", gotPatterns, tc.wantOutput)
			}
		})
	}
}

// TestCombinations tests combinations() function.
func TestCombinations(t *testing.T) {
	cases := []struct {
		name         string
		input        [][]string
		maxLiterals  int
		wantOutput   []string
		wantOverflow bool
	}{
		{
			name:        "1x1",
			input:       [][]string{{"A"}, {"B"}},
			maxLiterals: 100,
			wantOutput:  []string{"AB"},
		},
		{
			name:        "0x1",
			input:       [][]string{{}, {"B"}},
			maxLiterals: 100,
			wantOutput:  []string{},
		},
		{
			name:        "1x2",
			input:       [][]string{{"A"}, {"1", "2"}},
			maxLiterals: 100,
			wantOutput:  []string{"A1", "A2"},
		},
		{
			name:        "2x2",
			input:       [][]string{{"A", "B"}, {"1", "2"}},
			maxLiterals: 100,
			wantOutput:  []string{"A1", "A2", "B1", "B2"},
		},
		{
			name:        "empty string as an option",
			input:       [][]string{{"A", ""}, {"1", "2"}},
			maxLiterals: 100,
			wantOutput:  []string{"A1", "A2", "1", "2"},
		},
		{
			name:         "overflow",
			input:        [][]string{{"A", "B"}, {"1", "2"}},
			maxLiterals:  3,
			wantOverflow: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			gotPatterns, ok := combinations(tc.input, tc.maxLiterals)
			if tc.wantOverflow {
				if ok {
					t.Fatalf("expected to get an overflow, got success")
				}

				return
			}

			if !ok {
				t.Fatalf("unexpected overflow")
			}

			sort.Strings(tc.wantOutput)
			sort.Strings(gotPatterns)
			if !reflect.DeepEqual(tc.wantOutput, gotPatterns) {
				t.Fatalf("returned list of patterns (%#v) does not match the expected value (%#v)", gotPatterns, tc.wantOutput)
			}
		})
	}
}

// TestUnwrapCase tests unwrapping literals of case-insensitive regexps.
func TestUnwrapCase(t *testing.T) {
	cases := []struct {
		name          string
		ignoreCase    bool
		inputPatterns []string
		maxLiterals   int
		wantOutput    []string
		wantOverflow  bool
	}{
		{
			name:          "simple phrase",
			inputPatterns: []string{"simple phrase"},
			maxLiterals:   100,
			wantOutput:    []string{"simple phrase"},
		},
		{
			name:          "ignore case",
			ignoreCase:    true,
			inputPatterns: []string{"i"},
			maxLiterals:   100,
			wantOutput:    []string{"i", "I"},
		},
		{
			name:          "ignore case two letters",
			ignoreCase:    true,
			inputPatterns: []string{"ic"},
			maxLiterals:   100,
			wantOutput:    []string{"IC", "Ic", "iC", "ic"},
		},
		{
			name:          "ignore case two words",
			ignoreCase:    true,
			inputPatterns: []string{"i", "c"},
			maxLiterals:   100,
			wantOutput:    []string{"C", "I", "c", "i"},
		},
		{
			name:          "ignore case overflow",
			ignoreCase:    true,
			inputPatterns: []string{"long text"},
			maxLiterals:   100,
			wantOverflow:  true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			re := &syntax.Regexp{}
			if tc.ignoreCase {
				re.Flags = syntax.FoldCase
			}

			gotPatterns, ok := unwrapCase(re, tc.inputPatterns, tc.maxLiterals)
			if tc.wantOverflow {
				if ok {
					t.Fatalf("expected to get an overflow, got success")
				}

				return
			}

			if !ok {
				t.Fatalf("unexpected overflow")
			}

			sort.Strings(tc.wantOutput)
			sort.Strings(gotPatterns)
			if !reflect.DeepEqual(tc.wantOutput, gotPatterns) {
				t.Fatalf("returned list of patterns (%#v) does not match the expected value (%#v)", gotPatterns, tc.wantOutput)
			}
		})
	}
}

// TestFindLongestCommonLiteral tests finding longest literal in a regexp.
func TestFindLongestCommonLiteral(t *testing.T) {
	cases := []struct {
		input      string
		wantOutput string
	}{
		{
			input:      "simple phrase",
			wantOutput: "simple phrase",
		},
		{
			input:      "simple (phrase)?",
			wantOutput: "simple ",
		},
		{
			input:      "[iI]",
			wantOutput: "",
		},
		{
			input:      "[i]b",
			wantOutput: "ib",
		},
		{
			input:      "simple (phrase)+",
			wantOutput: "simple ",
		},
		{
			input:      "a*",
			wantOutput: "",
		},
		{
			input:      "(abc)|(ab)",
			wantOutput: "",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			re, err := syntax.Parse(tc.input, syntax.Perl)
			if err != nil {
				t.Fatalf("failed to parse regexp %q: %v", tc.input, err)
			}

			gotOutput := findLongestCommonLiteral(re)

			if gotOutput != tc.wantOutput {
				t.Fatalf("returned value (%q) does not match the expected value (%q)", gotOutput, tc.wantOutput)
			}
		})
	}
}

func contains(list []int, value int) bool {
	for _, elem := range list {
		if elem == value {
			return true
		}
	}
	return false
}

func TestPatterns(t *testing.T) {
	// Loading all crawlers with go:embed
	// some validation happens in UnmarshalJSON.
	allCrawlers := Crawlers

	// There are at least 10 crawlers.
	if len(allCrawlers) < 10 {
		t.Errorf("Number of crawlers must be at least 10, got %d.", len(allCrawlers))
	}

	if IsCrawler(browserUA) {
		t.Errorf("Browser UA %q was detected as a crawler.", browserUA)
	}
	if len(MatchingCrawlers(browserUA)) != 0 {
		t.Errorf("MatchingCrawlers found crawlers matching Browser UA %q.", browserUA)
	}

	for i, crawler := range allCrawlers {
		t.Run(crawler.Pattern, func(t *testing.T) {
			fmt.Println(crawler.Pattern)

			for _, instance := range crawler.Instances {
				if !IsCrawler(instance) {
					t.Errorf("Instance %q is not detected as a crawler.", instance)
				}
				hits := MatchingCrawlers(instance)
				if !contains(hits, i) {
					t.Errorf("Crawler with index %d (pattern %q) is not in the list returned by MatchingCrawlers(%q): %v.", i, crawler.Pattern, instance, hits)
				}
			}
		})
	}
}

func TestFalseNegatives(t *testing.T) {
	const browsersURL = "https://raw.githubusercontent.com/microlinkhq/top-user-agents/master/src/index.json"
	resp, err := http.Get(browsersURL)
	if err != nil {
		t.Fatalf("Failed to fetch the list of browser User Agents from %s: %v.", browsersURL, err)
	}

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	})

	var browsers []string
	if err := json.NewDecoder(resp.Body).Decode(&browsers); err != nil {
		t.Fatalf("Failed to parse the list of browser User Agents: %v.", err)
	}

	for _, userAgent := range browsers {
		if IsCrawler(userAgent) {
			t.Errorf("Browser User Agent %q is recognized as a crawler.", userAgent)
		}
		indices := MatchingCrawlers(userAgent)
		if len(indices) != 0 {
			t.Errorf("Browser User Agent %q matches with crawlers %v.", userAgent, indices)
		}
	}
}

const (
	crawlerUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/"
	browserUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) obsidian/1.5.3 Chrome/114.0.5735.289 Electron/25.8.1 Safari/537.36"
)

func BenchmarkIsCrawlerPositive(b *testing.B) {
	b.SetBytes(int64(len(crawlerUA)))
	for n := 0; n < b.N; n++ {
		if !IsCrawler(crawlerUA) {
			b.Fail()
		}
	}
}

func BenchmarkMatchingCrawlersPositive(b *testing.B) {
	b.SetBytes(int64(len(crawlerUA)))
	for n := 0; n < b.N; n++ {
		if len(MatchingCrawlers(crawlerUA)) == 0 {
			b.Fail()
		}
	}
}

func BenchmarkIsCrawlerNegative(b *testing.B) {
	b.SetBytes(int64(len(browserUA)))
	for n := 0; n < b.N; n++ {
		if IsCrawler(browserUA) {
			b.Fail()
		}
	}
}

func BenchmarkMatchingCrawlersNegative(b *testing.B) {
	b.SetBytes(int64(len(browserUA)))
	for n := 0; n < b.N; n++ {
		if len(MatchingCrawlers(browserUA)) != 0 {
			b.Fail()
		}
	}
}
