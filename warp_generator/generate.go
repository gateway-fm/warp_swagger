package warp_generator //nolint:all

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/models"
	"github.com/gateway-fm/warp_swagger/warp_generator/external_packages"
	"github.com/gateway-fm/warp_swagger/warp_generator/handlers"
	"github.com/gateway-fm/warp_swagger/warp_generator/middlewares"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

func Templates(config *config_warp.Warp,
	daily, requests []string, handlersModels *models.Handlers) ([]templater.ITemplate, error) {
	var templates templater.Templates
	external, err := external_packages.GenerateExternalModels(config, daily, requests)
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}
	var hndlrs []templater.ITemplate

	for i := range handlersModels.Operations {
		hndlr, err := handlers.GenerateHandlers(handlersModels.Operations[i], config, handlersModels.Operations[i].OperationsPath)
		if err != nil {
			return nil, fmt.Errorf("failed while collecting templates: %w", err)
		}
		hndlrs = append(hndlrs, hndlr)
	}
	mdwrs, err := middlewares.GenerateMdws()
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}
	hndlrs = append(hndlrs, external)
	hndlrs = append(hndlrs, mdwrs)
	templates = templater.GetAll(hndlrs...)

	return templates, nil
}
