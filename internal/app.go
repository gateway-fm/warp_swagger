package internal

import (
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_swagger"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/proto_parser"
	"github.com/gateway-fm/warp_swagger/swagger"
	"github.com/gateway-fm/warp_swagger/warp_generator"
	"github.com/gateway-fm/warp_swagger/warp_generator/external_packages"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

type App struct {
	warpCfg   *config_warp.Warp
	swagCfg   *config_swagger.SwaggerCfg
	templates []templater.ITemplate
	swag      swagger.ISwagger
	dummy     swagger.IDummy
}

func (app *App) WarpCfg() *config_warp.Warp {
	return app.warpCfg
}
func (app *App) SwagCfg() *config_swagger.SwaggerCfg {
	return app.swagCfg
}

func NewApplication() (app *App, err error) {
	app = &App{
		warpCfg: &config_warp.Warp{},
		swagCfg: &config_swagger.SwaggerCfg{},
	}
	return
}

func (app *App) CallDummy() error {
	d := swagger.NewDummy()

	err := d.Call("swagger/swagger.yaml")
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	app.dummy = d
	return nil
}
func (app *App) SwaggerGenerate(args []string) error {
	swag := swagger.NewSwagger(app.SwagCfg())

	if err := swag.Generate(args[0]); err != nil {
		return fmt.Errorf("failed to generate swagger:%w", err)
	}
	return nil
}

func ProtoMessageUnwrap(parsed []*proto_parser.ParsedMsg) ([]string, error) {
	var protoMessage []string
	var n, t []string
	for i := range parsed {
		n = append(n, parsed[i].ParsedNames)
		t = append(t, parsed[i].ParsedTypes)
	}
	protoMessage, err := external_packages.Merge(n, t)
	if err != nil {
		return nil, fmt.Errorf("failed to merge names and types:%w", err)
	}
	return protoMessage, nil
}

func (app *App) prepareTemplates() error {
	if err := app.CallDummy(); err != nil {
		return fmt.Errorf("failed to call dummy: %w", err)
	}
	protoParser := proto_parser.NewIProtoParser()
	requestModels, err := protoParser.Parse(app.WarpCfg().External.ProtoPath, "FetchedReqModelProto") // TODO   H A R D C O D E D ! ! !
	if err != nil {
		return fmt.Errorf("PrepareTemplates returned an error:%w", err)
	}

	dailyModels, err := protoParser.Parse(app.WarpCfg().External.ProtoPath, "DailyModelsProto")
	if err != nil {
		return fmt.Errorf("PrepareTemplates returned an error:%w", err)
	}

	daily, err := ProtoMessageUnwrap(dailyModels)
	if err != nil {
		return fmt.Errorf("failed to get daily proto message %w", err)
	}
	requests, err := ProtoMessageUnwrap(requestModels)
	if err != nil {
		return fmt.Errorf("failed to get daily proto message %w", err)
	}
	templates, err := warp_generator.Templates(app.WarpCfg(), daily, requests, app.dummy.GetHandlersModel())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	app.templates = templates
	// for i := range app.templates {
	// 	fmt.Println(app.templates[i])
	// }
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
