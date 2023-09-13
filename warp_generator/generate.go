package warp_generator //nolint:all

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/warp_generator/external_packages"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

func Templates(config *config_warp.Warp,
	simpleNames, simpleTypes, customNames, customTypes []string) ([]templater.ITemplate, error) {
	external, err := external_packages.GenerateExternalModels(config, simpleNames, simpleTypes, customNames, customTypes)
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}
	templates := templater.GetAll(external)
	return templates, nil
}
