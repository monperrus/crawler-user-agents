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

func BenchmarkIsCrawler(b *testing.B) {
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google-PageRenderer Google (+https://developers.google.com/+/web/snippet/)"
	b.SetBytes(int64(len(userAgent)))
	for n := 0; n < b.N; n++ {
		IsCrawler(userAgent)
	}
}

func BenchmarkMatchingCrawlers(b *testing.B) {
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google-PageRenderer Google (+https://developers.google.com/+/web/snippet/)"
	b.SetBytes(int64(len(userAgent)))
	for n := 0; n < b.N; n++ {
		MatchingCrawlers(userAgent)
	}
}
