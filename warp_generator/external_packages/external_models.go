package external_packages //nolint:all

import (
	"errors"
	"fmt"
	"github.com/gateway-fm/warp_swagger/config_warp"
	"github.com/gateway-fm/warp_swagger/proto_parser"
	"github.com/gateway-fm/warp_swagger/warp_generator/templater"
)

// ModelsFunc is
func ModelsFunc(m ...string) func() []string {
	foo := func() []string {
		return m
	}
	return foo
}
func Merge(names, types []string) ([]string, error) {
	pkgModels, err := mergeWithNames(names, types)
	if err != nil {
		return nil, fmt.Errorf("failed while merging: %w", err)
	}
	return pkgModels, nil
}

func GenerateExternalModels(
	config *config_warp.Warp,
	daily, requests []string,
) (templater.ITemplate, error) {
	path := "./templates/external_pkg_models.gohtml"

	var PkgNameUC = func() string {
		return config.External.PackageName
	}
	var ProtoName = func() string {
		return config.External.ProtoName
	}
	var ProtoPath = func() string {
		return config.External.ProtoPath
	}
	var PackageURL = func() string {
		return config.External.PackageURL
	}
	var Daily = func() []string {
		return daily
	}
	var Requests = func() []string {
		return requests
	}

	var funcNames = []string{
		"PkgNameUC",
		"ProtoName",
		"ProtoPath",
		"PackageURL",
		"Daily",
		"Requests",
	}

	funcs := templater.GetTemplateInterfaces(
		PkgNameUC,
		ProtoName,
		ProtoPath,
		PackageURL,
		Daily,
		Requests,
	)
	funcMap := templater.CompleteFuncMap(funcNames, funcs)
	elems := "external_pkg_models"
	output := fmt.Sprintf("pkg/%s/models/models.go", config.External.PackageName)
	ifaces := templater.GetTemplateInterfaces(daily, requests)
	template := templater.NewTemplate(path, output, ifaces, funcMap, elems)
	return template, nil
}

// TODO: add array check if keyword 'repeated' is presented
// protoConvertCases converts given proto-type into go-type
func protoConvertCases(s string) string {
	t := proto_parser.ParseProtoType(s)
	switch t {
	case proto_parser.ProtoBytes:
		return "[]byte"
	case proto_parser.ProtoInt64:
		//   WARNING! it could be time represented as timestamp
		//
		//   TODO: Add timestamp check!!!
		return "int"
	case proto_parser.ProtoFloat:
		return "float64"
	default:
		return s
	}
	//   todo: add error check lately
}

func convertProtoTypesToGoTypes(protoTypes []string) ([]string, error) {
	var converted []string
	for i := range protoTypes {
		converted = append(converted, protoConvertCases(protoTypes[i]))
	}
	if len(converted) != len(protoTypes) {
		return nil, ErrTypesLenMismatched
	}
	return converted, nil
}

func mergeWithNames(protoNames, protoTypes []string) ([]string, error) {
	var merged []string
	converted, err := convertProtoTypesToGoTypes(protoTypes)
	if err != nil {
		return nil, fmt.Errorf("merging names for external package models has been failed:%w", err)
	}

	//   converted len must be equal to protoNames
	for i := range converted {
		merged = append(merged, mergedName(protoNames[i], converted[i]))
	}

	return merged, nil
}
func mergedName(naMe, tyPe string) string {
	return fmt.Sprintf("%s 	%s", naMe, tyPe)
}

var ErrTypesLenMismatched = errors.New("proto and package types len are mismatched")
