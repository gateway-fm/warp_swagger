package main

import (
	"embed"
	"github.com/gateway-fm/warp_swagger/cmd/dummy"
	"github.com/gateway-fm/warp_swagger/cmd/root"
	"github.com/gateway-fm/warp_swagger/cmd/summon"
	"github.com/gateway-fm/warp_swagger/cmd/swagger"
	"github.com/gateway-fm/warp_swagger/internal"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
	"github.com/misnaged/annales/logger"
	"os"
)

//go:embed templates
var header embed.FS

func main() {
	app, err := internal.NewApplication()
	if err != nil {
		logger.Log().Errorf("An error occurred %v", err)
		os.Exit(1)
	}
	if err = templater.TempDir(); err != nil {
		logger.Log().Errorf("An error occurred when exec templater.TempDir %v", err)
		os.Exit(1)
	}

	rootCmd := root.Cmd(app)
	rootCmd.AddCommand(summon.Cmd(app, header))
	rootCmd.AddCommand(swagger.Cmd(app))
	rootCmd.AddCommand(dummy.Cmd(app))

	if err = rootCmd.Execute(); err != nil {
		logger.Log().Errorf("An error occurred %v", err)
		os.Exit(1)
	}
}
