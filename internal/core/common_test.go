package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/core"
)

func TestConvertMentionsToNames(t *testing.T) {
	input := `<@1020> <@!1030> <@2030> <@2040>`

	assert.Equal(t, input, core.ConvertMentionsToNames(input, func(mention string) string { return mention }))
	assert.Equal(t, "Peter Klaus Theo Ferdinand", core.ConvertMentionsToNames(input,
		func(mention string) string {
			switch mention {
			case "1020":
				return "Peter"
			case "1030":
				return "Klaus"
			case "2030":
				return "Theo"
			case "2040":
				return "Ferdinand"
			default:
				return mention
			}
		}))
}
