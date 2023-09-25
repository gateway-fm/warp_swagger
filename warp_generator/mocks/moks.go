package mocks

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
	"github.com/go-openapi/inflect"
	"strings"
)

func GenerateMocks(o string, mocks []string, config *config_warp.Warp) (templater.ITemplate, error) {
	path := "templates/mocks.gohtml"
	var Mocks = func() []string {
		return mocks
	}
	var PackageURL = func() string {
		return config.External.PackageURL
	}
	var Model = func() string {
		spl := strings.Split(o, "_")
		for i := range spl {
			if spl[i] == "_" || spl[i] == "-" {
				spl = spl[i:]
			}
			spl[i] = inflect.Capitalize(spl[i])

		}
		o = strings.Join(spl, "")
		return o
	}
	var funcNames = []string{
		"Mocks",
		"Model",
		"PackageURL",
	}
	funcs := templater.GetTemplateInterfaces(
		Mocks,
		Model,
		PackageURL,
	)
	output := fmt.Sprintf("internal/mocks/%s.go", o)
	funcMap := templater.CompleteFuncMap(funcNames, funcs)
	elems := "mocks_main"
	ifaces := templater.GetTemplateInterfaces(elems)
	template := templater.NewTemplate(path, output, ifaces, funcMap, elems)
	return template, nil
}
