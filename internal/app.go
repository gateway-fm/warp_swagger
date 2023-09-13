package internal

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/proto_parser"
	"github.com/gateway-fm/warp_swagger/warp_generator"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

type App struct {
	warpCfg   *config_warp.Warp
	templates []templater.ITemplate
}

func (app *App) WarpCfg() *config_warp.Warp {
	return app.warpCfg
}

func NewApplication() (app *App, err error) {
	app = &App{
		warpCfg: &config_warp.Warp{},
	}
	return
}

func (app *App) prepareTemplates() error {
	protoParser := proto_parser.NewIProtoParser()
	fmt.Println(app.WarpCfg().External.ProtoPath)

	primitive, custom, err := protoParser.Parse(app.WarpCfg().External.ProtoPath)
	if err != nil {
		return fmt.Errorf("PrepareTemplates returned an error:%w", err)
	}

	//  TODO: refactor: get rid off ***simpleNames, simpleTypes, customNames, customTypes []string***
	//   using something  similar to []*ParsedMsg is preferable
	var simpleNames, simpleTypes, customNames, customTypes []string

	for i := range primitive {
		simpleNames = append(simpleNames, primitive[i].ParsedNames)
		simpleTypes = append(simpleTypes, primitive[i].ParsedTypes)
	}
	for i := range custom {
		customNames = append(customNames, custom[i].ParsedNames)
		customTypes = append(customTypes, custom[i].ParsedTypes)
	}

	//   *******************-- TO DO --********************* //

	templates, err := warp_generator.Templates(app.WarpCfg(), simpleNames, simpleTypes, customNames, customTypes)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	app.templates = templates
	return nil
}

func (app *App) Summon() error {
	if err := app.prepareTemplates(); err != nil {
		return fmt.Errorf("summon was interrupted:%w", err)
	}
	for _, t := range app.templates {
		err := t.Generate()
		if err != nil {
			return fmt.Errorf("summon was interrupted:%w", err)
		}
	}
	return nil
}
