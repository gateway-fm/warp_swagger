package handlers

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/models"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
	"github.com/go-openapi/inflect"
)

func GenerateHandlers(
	operation *models.Operation,
	config *config_warp.Warp,
	operationPath string,
) (templater.ITemplate, error) {
	path := "templates/internal_handlers.gohtml"

	var OperationIDuc = func() string {
		return inflect.Capitalize(operation.OperationID)
	}
	var OperationIDlc = func() string {
		return operation.OperationID
	}
	var PackageURL = func() string {
		return config.External.PackageURL
	}
	var OperationsPath = func() string {
		return operationPath
	}
	var funcNames = []string{
		"OperationIDuc",
		"OperationIDlc",
		"PackageURL",
		"OperationsPath",
	}

	funcs := templater.GetTemplateInterfaces(
		OperationIDuc,
		OperationIDlc,
		PackageURL,
		OperationsPath,
	)
	funcMap := templater.CompleteFuncMap(funcNames, funcs)
	elems := "handler_main"
	ifaces := templater.GetTemplateInterfaces(operation)
	output := fmt.Sprintf("internal/handlers/%s", operation.OutputFileName)
	template := templater.NewTemplate(path, output, ifaces, funcMap, elems)
	return template, nil
}
