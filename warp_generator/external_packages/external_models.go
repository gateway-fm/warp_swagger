package external_packages //nolint:revive

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

func GenerateExternalModels(
	config *config_warp.Warp,
	simpleNames, simpleTypes, customNames, customTypes []string,
) (templater.ITemplate, error) {
	path := "templates/pkg/models/external_pkg_models.gohtml"

	ProtosModelStruct := struct {
		SimpleNames, SimpleTypes, CustomNames, CustomTypes []string
	}{
		SimpleNames: simpleNames,
		SimpleTypes: simpleTypes,
		CustomNames: customNames,
		CustomTypes: customTypes,
	}
	//   TODO: handle custom Names&Types!

	pkgModels, err := mergeWithNames(ProtosModelStruct.SimpleNames, ProtosModelStruct.SimpleTypes)
	if err != nil {
		return nil, fmt.Errorf("failed while merging: %w", err)
	}
	var PkgNameUC = func() string {
		return config.External.PackageName
	}
	var ProtoName = func() string {
		return config.External.ProtoName
	}
	var PackageURL = func() string {
		return config.External.PackageURL
	}
	var PkgModels = func() []string {
		return pkgModels
	}
	for i := range pkgModels {
		fmt.Println(pkgModels[i])
	}
	var funcNames = []string{
		"PkgNameUC",
		"ProtoName",
		"PackageURL",
		"PkgModels",
	}

	funcs := templater.GetTemplateInterfaces(
		PkgNameUC,
		ProtoName,
		PackageURL,
		PkgModels,
	)
	funcMap := templater.CompleteFuncMap(funcNames, funcs)
	elems := "external_pkg_models"
	ifaces := templater.GetTemplateInterfaces(ProtosModelStruct)
	template := templater.NewTemplate(path, config.External.Output, ifaces, funcMap, elems)
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
