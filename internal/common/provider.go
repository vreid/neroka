package common

import (
	"log"

	"github.com/samber/lo"
)

type BaseProvider struct {
	AvailableModels map[string]bool
	DefaultModel    string
}

func (p *BaseProvider) AvaiableModels() []string {
	return lo.Keys(p.AvailableModels)
}

func (p *BaseProvider) CheckDefaultModel() bool {
	if _, ok := p.AvailableModels[p.DefaultModel]; !ok {
		log.Printf("coudln't find default model '%s' in available models\n", p.DefaultModel)
		return false
	}

	return true
}
