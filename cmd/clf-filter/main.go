// clf-filter reads Combined Log Format lines from stdin and writes them to stdout,
// removing bot/crawler lines by default. Use --bot to keep only bot lines.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	agents "github.com/monperrus/crawler-user-agents"
)

func extractUserAgent(line string) (string, bool) {
	// Combined Log Format ends with: "referer" "user-agent"
	// Find the last quoted field.
	end := strings.LastIndex(line, "\"")
	if end < 1 {
		return "", false
	}
	start := strings.LastIndex(line[:end], "\"")
	if start < 0 {
		return "", false
	}
	return line[start+1 : end], true
}

func main() {
	botOnly := flag.Bool("bot", false, "keep only bot/crawler lines (default: remove bots)")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	// Support long lines (e.g. large URLs).
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		ua, ok := extractUserAgent(line)
		isBot := ok && agents.IsCrawler(ua)

		if *botOnly == isBot {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "clf-filter: read error:", err)
		os.Exit(1)
	}
}
