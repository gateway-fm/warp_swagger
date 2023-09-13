package config_warp // nolint:all

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("external.output", "models.go")
	viper.SetDefault("external.packagename", "observerApi")
	viper.SetDefault("external.protoname", "observer")
	viper.SetDefault("external.protopath", "model.proto")
	viper.SetDefault("external.packageurl", "github.com/gateway-fm/warp_swagger")

}
