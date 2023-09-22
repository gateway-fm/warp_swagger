package middlewares

import (
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

func GenerateMdws() (templater.ITemplate, error) {
	path := "templates/internal_middlewares.gohtml"

	elems := "middlewares_main"
	ifaces := templater.GetTemplateInterfaces(elems)
	output := "internal/middlewares/middlewares.go"
	template := templater.NewTemplate(path, output, ifaces, nil, elems)
	return template, nil
}
