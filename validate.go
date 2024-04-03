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

func joinRes(begin, end int) string {
	regexps := make([]string, 0, len(Crawlers))
	for _, crawler := range Crawlers[begin:end] {
		regexps = append(regexps, "("+crawler.Pattern+")")
	}
	return strings.Join(regexps, "|")
}

var allRegexps = joinRes(0, len(Crawlers))

var allRegexpsRe = regexp.MustCompile(allRegexps)

// Returns if User Agent string matches any of crawler patterns.
func IsCrawler(userAgent string) bool {
	return allRegexpsRe.MatchString(userAgent)
}

// With RE2 it is fast to check the text against a large regexp.
// To find matching regexps faster, built a binary tree of regexps.

type regexpNode struct {
	re    *regexp.Regexp
	left  *regexpNode
	right *regexpNode
	index int
}

var regexpsTree = func() *regexpNode {
	nodes := make([]*regexpNode, len(Crawlers))
	starts := make([]int, len(Crawlers)+1)
	for i, crawler := range Crawlers {
		nodes[i] = &regexpNode{
			re:    regexp.MustCompile(crawler.Pattern),
			index: i,
		}
		starts[i] = i
	}
	starts[len(Crawlers)] = len(Crawlers) // To get end of interval.

	for len(nodes) > 1 {
		// Join into pairs.
		nodes2 := make([]*regexpNode, (len(nodes)+1)/2)
		starts2 := make([]int, 0, len(nodes2)+1)
		for i := 0; i < len(nodes)/2; i++ {
			leftIndex := 2 * i
			rightIndex := 2*i + 1
			nodes2[i] = &regexpNode{
				left:  nodes[leftIndex],
				right: nodes[rightIndex],
			}
			if len(nodes2) != 1 {
				// Skip regexp for root node, it is not used.
				joinedRe := joinRes(starts[leftIndex], starts[rightIndex+1])
				nodes2[i].re = regexp.MustCompile(joinedRe)
			}
			starts2 = append(starts2, starts[leftIndex])
		}
		if len(nodes)%2 == 1 {
			nodes2[len(nodes2)-1] = nodes[len(nodes)-1]
			starts2 = append(starts2, starts[len(starts)-2])
		}
		starts2 = append(starts2, starts[len(starts)-1])

		nodes = nodes2
		starts = starts2
	}

	root := nodes[0]

	if root.left == nil {
		panic("the algoriths does not work with just one regexp")
	}

	return root
}()

// Finds all crawlers matching the User Agent and returns the list of their indices in Crawlers.
func MatchingCrawlers(userAgent string) []int {
	indices := []int{}

	var visit func(node *regexpNode)
	visit = func(node *regexpNode) {
		if node.left != nil {
			if node.left.re.MatchString(userAgent) {
				visit(node.left)
			}
			if node.right.re.MatchString(userAgent) {
				visit(node.right)
			}
		} else {
			// Leaf.
			indices = append(indices, node.index)
		}
	}

	visit(regexpsTree)

	return indices
}
