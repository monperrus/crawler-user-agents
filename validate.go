package agents

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	regexp "github.com/wasilibs/go-re2"
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

var allRegexps = func() string {
	regexps := make([]string, 0, len(Crawlers))
	for _, crawler := range Crawlers {
		regexps = append(regexps, "("+crawler.Pattern+")")
	}
	return strings.Join(regexps, "|")
}()

var allRegexpsRe = regexp.MustCompile(allRegexps)

// Returns if User Agent string matches any of crawler patterns.
func IsCrawler(userAgent string) bool {
	return allRegexpsRe.MatchString(userAgent)
}

var individualRegexps = func() []*regexp.Regexp {
	regexps := make([]*regexp.Regexp, len(Crawlers))
	for i, crawler := range Crawlers {
		regexps[i] = regexp.MustCompile(crawler.Pattern)
	}
	return regexps
}()

// Finds all crawlers matching the User Agent and returns the list of their indices in Crawlers.
func MatchingCrawlers(userAgent string) []int {
	indices := []int{}
	for i, regexp := range individualRegexps {
		if regexp.MatchString(userAgent) {
			indices = append(indices, i)
		}
	}

	return indices
}
