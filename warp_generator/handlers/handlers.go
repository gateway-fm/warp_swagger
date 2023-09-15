package handlers

import (
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
	path := "templates/internal/handler_template.gohtml"

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
	template := templater.NewTemplate(path, operation.OutputFileName, ifaces, funcMap, elems)
	return template, nil
}
