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
