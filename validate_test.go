package agents

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatterns(t *testing.T) {
	// Loading all crawlers with go:embed
	// some validation happens in UnmarshalJSON.
	allCrawlers := Crawlers

	// There are at least 10 crawlers.
	require.GreaterOrEqual(t, len(allCrawlers), 10)

	for i, crawler := range allCrawlers {
		t.Run(crawler.Pattern, func(t *testing.T) {
			fmt.Println(crawler.Pattern)

			for _, instance := range crawler.Instances {
				require.True(t, IsCrawler(instance), instance)
				require.Contains(t, MatchingCrawlers(instance), i, instance)
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
