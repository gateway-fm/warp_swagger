package config_swagger //nolint:revive

import "github.com/go-openapi/spec"

type SwaggerCfg struct {
	Spec *spec.Swagger
}

func (s *SwaggerCfg) GetSpec() {
}
