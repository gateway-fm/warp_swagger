package config_swagger //nolint:all

import "github.com/go-openapi/spec"

type SwaggerCfg struct {
	Spec *spec.Swagger
}

func (s *SwaggerCfg) GetSpec() {
}
