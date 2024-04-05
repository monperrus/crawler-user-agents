package agents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

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
