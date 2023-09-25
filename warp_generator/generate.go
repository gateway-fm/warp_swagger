package warp_generator //nolint:all

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/models"
	app_gen "github.com/gateway-fm/warp_swagger/warp_generator/app-gen"
	"github.com/gateway-fm/warp_swagger/warp_generator/external_packages"
	"github.com/gateway-fm/warp_swagger/warp_generator/handlers"
	"github.com/gateway-fm/warp_swagger/warp_generator/middlewares"
	"github.com/gateway-fm/warp_swagger/warp_generator/mocks"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

func Templates(config *config_warp.Warp,
	daily, requests []string, handlersModels *models.Handlers, api *models.API, m map[string][]string) ([]templater.ITemplate, error) {
	var templates templater.Templates
	external, err := external_packages.GenerateExternalModels(config, daily, requests)
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}
	var hndlrs []templater.ITemplate
	var operIds []string
	for i := range handlersModels.Operations {
		hndlr, err := handlers.GenerateHandlers(handlersModels.Operations[i], config, handlersModels.Operations[i].OperationsPath)
		if err != nil {
			return nil, fmt.Errorf("failed while collecting templates: %w", err)
		}
		operIds = append(operIds, handlersModels.Operations[i].OperationID)
		hndlrs = append(hndlrs, hndlr)
	}
	appGen, err := app_gen.GenerateAppHandlers(api, config, operIds)
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}
	mdwrs, err := middlewares.GenerateMdws()
	if err != nil {
		return nil, fmt.Errorf("failed while collecting templates: %w", err)
	}

	for model, modelMocks := range m {
		moks, err := mocks.GenerateMocks(model, modelMocks, config)
		if err != nil {
			return nil, fmt.Errorf("failed while collecting templates: %w", err)
		}
		hndlrs = append(hndlrs, moks)
	}
	hndlrs = append(hndlrs, appGen)
	hndlrs = append(hndlrs, external)
	hndlrs = append(hndlrs, mdwrs)
	templates = templater.GetAll(hndlrs...)

	return templates, nil
}
