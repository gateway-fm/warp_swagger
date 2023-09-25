package app_gen //nolint:all
import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/models"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
	"github.com/go-openapi/inflect"
)

func GenerateAppHandlers(
	api *models.API,
	config *config_warp.Warp,
	operationIDs []string,
) (templater.ITemplate, error) {
	path := "templates/internal_apphandlers.gohtml"

	var OperationIDs = func() []string {
		var capitalized []string
		for i := range operationIDs {
			capitalized = append(capitalized, inflect.Capitalize(operationIDs[i]))
		}
		return capitalized
	}
	var PackageURL = func() string {
		return config.External.PackageURL
	}

	var NameAPI = func() string {
		return api.Name
	}
	var funcNames = []string{
		"OperationIDs",
		"NameAPI",
		"PackageURL",
	}

	funcs := templater.GetTemplateInterfaces(
		OperationIDs,
		NameAPI,
		PackageURL,
	)
	funcMap := templater.CompleteFuncMap(funcNames, funcs)
	elems := "app_handlers_main"
	ifaces := templater.GetTemplateInterfaces(api)
	output := fmt.Sprintf("internal/%s", "app_handlers.go")
	template := templater.NewTemplate(path, output, ifaces, funcMap, elems)
	return template, nil
}
