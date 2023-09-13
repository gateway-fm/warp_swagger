package utils

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
	b, _ := os.ReadFile(path)
	return b
}

func GoFiles(pathPart string) []string {
	cfg := &packages.Config{Mode: packages.NeedFiles}
	var Goes []string
	path := fmt.Sprintf("./%s/models", pathPart)
	pkgs, err := packages.Load(cfg, path)
	if err == nil {
		for _, pkg := range pkgs {
			Goes = append(Goes, pkg.GoFiles...)
		}
	}
	return Goes
}
func GenMod(cfg, generatePath string) *generate.Model {
	model := &generate.Model{}
	_, _ = flags.Parse(model)
	model.Shared.Target = flags.Filename(generatePath)
	model.Shared.Spec = flags.Filename(cfg)
	//  opts := new(generator.GenOpts)
	//  fmt.Println("spec path is", (*opts).Models)
	//  err := generator.GenerateModels(model.Name, new())
	//  if err != nil {
	//  	log.Fatalln(err)
	//  }

	if err := model.Execute([]string{}); err != nil {
		logger.Log().Errorf("error %v", err)
	}

	return model
}

func UnwrapAst(file *ast.File) {
	//   todo: comment this madness

	for i := range file.Decls {

		d := file.Decls[i]
		switch d.(type) { //nolint:gocritic
		case *ast.FuncDecl: //nolint:gocritic
		case *ast.GenDecl: //nolint:gocritic
			dd := d.(*ast.GenDecl).Specs
			for ii := range dd {

				spc := dd[ii]

				switch spc.(type) { //nolint:gocritic
				case *ast.ImportSpec: //nolint:gocritic
				case *ast.ValueSpec: //nolint:gocritic
				case *ast.TypeSpec: //nolint:gocritic
					tp := spc.(*ast.TypeSpec).Type
					list := tp.(*ast.StructType).Fields.List
					for iii := range list {
						for namesIdx := range list[iii].Names {
							expression := list[iii].Type
							switch expression.(type) {
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
								arrString = fmt.Sprintf("[]%s", arrString)
								fmt.Println(list[iii].Names[namesIdx].String(), arrString)
							case *ast.MapType:
								fmt.Println(list[iii].Names[namesIdx].String(), expression.(*ast.MapType).Key, expression.(*ast.MapType).Value) //nolint:gosimple
							case *ast.SliceExpr:
								fmt.Println(list[iii].Names[namesIdx].String(), expression.(*ast.SliceExpr).X) //nolint:gosimple
							case *ast.StarExpr:
								starPart := expression.(*ast.StarExpr).X  //nolint:gosimple
								starString := fmt.Sprintf("%s", starPart) //   we can't use starPart as a string
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

								fmt.Println(list[iii].Names[namesIdx].String(), str)
							}
						}
					}
				}
			}
		}
	}

}
