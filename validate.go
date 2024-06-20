package agents

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// Private time needed to convert addition_date from/to the format used in JSON.
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

var pattern2literals = map[string][]string{
	`[wW]get`:               {`wget`, `Wget`},
	`Ahrefs(Bot|SiteAudit)`: {`AhrefsBot`, `AhrefsSiteAudit`},
	`S[eE][mM]rushBot`:      {`SemrushBot`, `SeMrushBot`, `SEmrushBot`, `SEMrushBot`},
	`Livelap[bB]ot`:         {`Livelapbot`, `LivelapBot`},
	`[pP]ingdom`:            {`pingdom`, `Pingdom`},
	`Bark[rR]owler`:         {`Barkrowler`, `BarkRowler`},
	`^Apache-HttpClient`:    {`^Apache-HttpClient`},
	`^LCC `:                 {`^LCC `},
	`(^| )sentry\/`:         {`^sentry/`, ` sentry/`},
	`^curl`:                 {`^curl`},
	`[Cc]urebot`:            {`Curebot`, `curebot`},
	`^PHP-Curl-Class`:       {`^PHP-Curl-Class`},
	`(^| )PTST\/`:           {`^PTST/`, ` PTST/`},
	`^BW\/`:                 {`^BW/`},
}

var pattern2mainLiteral = map[string]string{
	`AdsBot-Google([^-]|$)`:             `AdsBot-Google`,
	`BlogTraffic\/\d\.\d+ Feed-Fetcher`: `BlogTraffic/`,
}

func analyzePattern(pattern string) (olds []string, re *regexp.Regexp) {
	literals, has := pattern2literals[pattern]
	if has {
		return literals, nil
	}

	re = regexp.MustCompile(pattern)
	prefix, complete := re.LiteralPrefix()
	if complete {
		return []string{prefix}, nil
	}

	mainLiteral, has := pattern2mainLiteral[pattern]
	if !has {
		panic("don't know what to do with pattern: " + pattern)
	}
	return []string{mainLiteral}, re
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
		literals, re := analyzePattern(crawler.Pattern)

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
