package templater

import (
	"fmt"
	"github.com/go-swagger/go-swagger/generator"
	"github.com/misnaged/annales/logger"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type ITemplate interface {
	Generate() error
	Excluded() bool
	GenerateNonGo() error
}

type Template struct {
	Elems                              []string
	Ifaces                             []interface{}
	ConfigTemplatePath, OutPutFilePath string
	FuncMap                            template.FuncMap
	IsExcluded                         bool
}

func ExecHeaderTemplate(output io.Writer) error {
	b, err := os.ReadFile("templates/header.gohtml")
	if err != nil {
		return fmt.Errorf("failed to read file:%w", err)
	}
	funcMap := generator.FuncMapFunc(generator.DefaultLanguageFunc())
	t := template.Must(template.New("").Funcs(funcMap).Parse(string(b)))

	if err = t.ExecuteTemplate(output, "", nil); err != nil {
		return fmt.Errorf("failed to execute template %w", err)
	}
	return nil
}
func (t *Template) GenerateFile() error {

	elems := t.Elems
	ifaces := t.Ifaces
	//   as an empty .go file and just "filled up" in this func
	file, _ := os.Create(t.OutPutFilePath) //
	defer file.Close()

	// if err := ExecHeaderTemplate(file); err != nil {
	// 	return fmt.Errorf("failed to exec ExecHeaderTemplate:%w", err)
	// }
	//    path to template file is absolute here, but it doesn't make any sense :D
	pattern, _ := filepath.Abs(t.ConfigTemplatePath) //   .gotmpl is used because of IDE's supports only :D

	//   template final preparation. Template must parse given pattern (which is our scheme.gotmpl file)
	tmpl := template.Must(template.New("").Funcs(t.FuncMap).ParseFiles(pattern))
	var wg sync.WaitGroup
	for i := range elems {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := tmpl.ExecuteTemplate(file, elems[i], ifaces[i]) //   first arg is output, second is the data we want to pass to this config. It could also be nil.
			if err != nil {
				logger.Log().Errorf("An error occurred %s", err)
				return
			}
		}()
		wg.Wait()
	}

	return nil
}
func (t *Template) Excluded() bool {
	return t.IsExcluded
}

// GoFmt execute go fmt for specific .go file from the code
func GoFmt(path string) error {

	read, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmted, err := format.Source(read)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, fmted, 0666) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

// Generate is main generation function
func (t *Template) Generate() error {
	err := t.GenerateFile()
	if err != nil {
		return fmt.Errorf(" Generate returned an error: %w", err)
	}
	err = GoFmt(t.OutPutFilePath)
	if err != nil {
		return fmt.Errorf(" go fmt returned an error: %w", err)
	}

	return nil
}

// GenerateNonGo is main generation function for non .go files
func (t *Template) GenerateNonGo() error {
	err := t.GenerateFile()
	if err != nil {
		return fmt.Errorf(" GenerateNonGo returned an error: %w", err)
	}
	return nil
}

// NewTemplate is
func NewTemplate(path, output string, ifaces []any, funcMap template.FuncMap, elems ...string) ITemplate {
	t := &Template{}
	t.ConfigTemplatePath = path
	t.OutPutFilePath = output
	t.Elems = elems
	t.Ifaces = ifaces
	t.FuncMap = funcMap
	return t
}

type Templates []ITemplate

func GetAll(templates ...ITemplate) Templates {
	templates = append(templates, templates...)
	return templates
}
