package core

import (
	"regexp"
)

var (
	mentionRegex = regexp.MustCompile(`<@!?(\d+)>`)
)

func ConvertMentionsToNames(text string, lookup func(string) string) string {
	if len(text) == 0 {
		return text
	}

	replaceMention := func(mention string) string {
		submatch := mentionRegex.FindStringSubmatch(mention)
		if len(submatch) < 2 {
			return mention
		}

		replacement := lookup(submatch[1])
		if replacement == submatch[1] {
			return mention
		} else {
			return replacement
		}
	}

	return mentionRegex.ReplaceAllStringFunc(text, replaceMention)
}
