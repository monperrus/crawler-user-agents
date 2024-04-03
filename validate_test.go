package agents

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatterns(t *testing.T) {
	for i, crawler := range Crawlers {
		t.Run(crawler.URL, func(t *testing.T) {
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
