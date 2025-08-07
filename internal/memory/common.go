package memory

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var (
	stopWords = []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for",
		"of", "with", "by", "is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did", "will", "would", "could",
		"should", "may", "might", "can", "shall", "must", "i", "you", "he",
		"she", "it", "we", "they", "me", "him", "her", "us", "them", "my",
		"your", "his", "her", "its", "our", "their", "this", "that", "these",
		"those", "what", "where", "when", "why", "how", "who", "which", "said",
		"says", "just", "like", "get", "got", "go", "went", "come", "came",
		"see", "saw", "know", "knew", "think", "thought", "take", "took",
		"make", "made", "give", "gave", "tell", "told", "ask", "asked",
		"discord", "server", "channel", "conversation", "summary", "between",
		"users", "bots",
	}

	cleanRegex = regexp.MustCompile(`[^\w\s]`)
)

func ExtractKeywords(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	cleanedText := cleanRegex.ReplaceAllString(strings.ToLower(text), " ")
	words := strings.Split(cleanedText, " ")

	keywords := lo.Filter(words, func(word string, _ int) bool {
		if len(word) < 3 {
			return false
		}

		if lo.Contains(stopWords, word) {
			return false
		}

		if _, err := strconv.Atoi(word); err == nil {
			return false
		}

		return true
	})

	keywords = lo.Uniq(keywords)

	keywordCount := int(math.Min(20.0, float64(len(keywords)))) // ugh
	return keywords[:keywordCount]
}
