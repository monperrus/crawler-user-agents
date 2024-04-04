package agents

import (
	"fmt"
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

const (
	crawlerUA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/"
	browserUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) obsidian/1.5.3 Chrome/114.0.5735.289 Electron/25.8.1 Safari/537.36"
)

func BenchmarkIsCrawlerPositive(b *testing.B) {
	b.SetBytes(int64(len(crawlerUA)))
	for n := 0; n < b.N; n++ {
		IsCrawler(crawlerUA)
	}
}

func BenchmarkMatchingCrawlersPositive(b *testing.B) {
	b.SetBytes(int64(len(crawlerUA)))
	for n := 0; n < b.N; n++ {
		MatchingCrawlers(crawlerUA)
	}
}

func BenchmarkIsCrawlerNegative(b *testing.B) {
	b.SetBytes(int64(len(browserUA)))
	for n := 0; n < b.N; n++ {
		IsCrawler(browserUA)
	}
}

func BenchmarkMatchingCrawlersNegative(b *testing.B) {
	b.SetBytes(int64(len(browserUA)))
	for n := 0; n < b.N; n++ {
		MatchingCrawlers(browserUA)
	}
}
