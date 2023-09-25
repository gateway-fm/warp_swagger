package swagger

import (
	"fmt"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	"github.com/jessevdk/go-flags"
	"github.com/misnaged/annales/logger"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
)

func BytesFromFile(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		logger.Log().Errorf("error while reading file:%v", err)
		return nil
	}
	return b
}

func GoFiles() []string {
	cfg := &packages.Config{Mode: packages.NeedFiles}
	var Goes []string
	path := fmt.Sprintf("./%s/models", "internal")
	pkgs, err := packages.Load(cfg, path)
	if err == nil {
		for _, pkg := range pkgs {
			Goes = append(Goes, pkg.GoFiles...)
		}
	}
	return Goes
}
func GenMod(cfg, generatePath string) (*generate.Model, error) {
	if cfg == "" || generatePath == "" {
		return nil, ErrNilGenerationPathOrOutput
	}
	model := &generate.Model{}
	_, err := flags.Parse(model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags %w", err)
	}
	model.Shared.Target = flags.Filename(generatePath)
	model.Shared.Spec = flags.Filename(cfg)
	model.Models.ModelPackage = "internal/models"
	if err = model.Execute([]string{}); err != nil {
		return nil, fmt.Errorf("executing model generation failed:%w", err)
	}

	return model, nil
}

func GenServer(cfg, generatePath string) (*generate.Server, error) {
	if cfg == "" || generatePath == "" {
		return nil, ErrNilGenerationPathOrOutput
	}
	server := &generate.Server{}
	_, err := flags.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags %w", err)
	}
	server.Shared.Target = flags.Filename(generatePath)
	server.Shared.Spec = flags.Filename(cfg)
	// server.Name
	server.ServerPackage = "internal/server"
	server.Models.ModelPackage = "internal/models"
	server.ExcludeMain = true
	if err = server.Execute([]string{}); err != nil {
		return nil, fmt.Errorf("executing server generation failed:%w", err)
	}

	return server, nil
}

func unwrapAst(file *ast.File) []string {
	var names []string
	for i := range file.Decls {

		d := file.Decls[i]
		switch d.(type) { //nolint:all
		case *ast.FuncDecl: //nolint:gocritic
		case *ast.GenDecl: //nolint:gocritic
			dd := d.(*ast.GenDecl).Specs
			for ii := range dd {

				spc := dd[ii]

				switch spc.(type) { //nolint:all
				case *ast.ImportSpec: //nolint:gocritic
				case *ast.ValueSpec: //nolint:gocritic
				case *ast.TypeSpec: //nolint:gocritic
					tp := spc.(*ast.TypeSpec).Type
					list := tp.(*ast.StructType).Fields.List
					for iii := range list {
						for namesIdx := range list[iii].Names {
							names = append(names, list[iii].Names[namesIdx].String())
							expression := list[iii].Type
							switch expression.(type) { //nolint:all
							case *ast.ArrayType:
								arrPart := expression.(*ast.ArrayType).Elt //nolint:gosimple
								arrString := fmt.Sprintf("%v", arrPart)
								if arrString[0] == '&' {
									arrStrings := strings.Split(arrString, " ")
									arrString = arrStrings[1]
									arrBytes := []byte(arrString)
									arrBytes = arrBytes[:len(arrBytes)-1]
									arrString = "*" + string(arrBytes)
								}
								arrString = fmt.Sprintf("[]%s", arrString) //nolint:all
								// fmt.Println(list[iii].Names[namesIdx].String(), arrString)
							case *ast.MapType:
								// fmt.Println(list[iii].Names[namesIdx].String(), expression.(*ast.MapType).Key, expression.(*ast.MapType).Value) //nolint:gosimple
							case *ast.SliceExpr:
								// fmt.Println(list[iii].Names[namesIdx].String(), expression.(*ast.SliceExpr).X) //nolint:gosimple
							case *ast.StarExpr:
								starPart := expression.(*ast.StarExpr).X  //nolint:gosimple
								starString := fmt.Sprintf("%s", starPart) // we can't use starPart as a string
								runed := []byte(starString)
								var newRuned []byte
								if string(runed[0]) == "&" {
									for runeIDX := range runed {
										if string(runed[runeIDX]) != "&" &&
											string(runed[runeIDX]) != "{" &&
											string(runed[runeIDX]) != "}" {
											newRuned = append(newRuned, runed[runeIDX])
										}
									}
									toMerge := strings.Split(string(newRuned), " ")
									starString = fmt.Sprintf("%s.%s", toMerge[0], toMerge[1])
								}
								str := fmt.Sprintf("*%s", starString)
								if str == "" { //nolint:all
								}
								// fmt.Println(list[iii].Names[namesIdx].String(), str)
							}
						}
					}
				}
			}
		}
	}
	return names
}
func RemoveDupes(sliceToCheck []string) []string {
	m := make(map[string]string)
	var cleanSlice []string
	for i := range sliceToCheck {
		m[sliceToCheck[i]] = ""
	}
	for i := range m {
		cleanSlice = append(cleanSlice, i)
	}
	return cleanSlice
}
