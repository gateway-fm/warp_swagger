/*



	 	Deprecated,

		but could be used in future
		if we would be needed to parse raw .yml

*/

package yaml_parser //nolint:revive

/*
import (
	"fmt"
	"gopkg.in/yaml.v3"
	"warp_swagger/config_swagger"
	"warp_swagger/models"
	"warp_swagger/utils"
	"warp_swagger/warp_errors"
)

type Parser struct {
	SwaggerCfg *config_swagger.SwaggerCfg
}

func NewParser(path string) (*Parser, error) {
	p := &Parser{}
	var err error
	p.SwaggerCfg, err = config_swagger.NewSwaggerCfg(path)
	if err != nil {
		return nil, fmt.Errorf("new swagger error: %w", err)
	}

	return p, nil
}

//   CollectRESTmethods searching any (post and get, at this time) kind of REST methods
//   in the 'paths' section of config (where they normally should be defined)
func (p *Parser) CollectRESTmethods() ([]map[string]any, error) {
	m := p.SwaggerCfg.SwagBP.PathsMap
	mapArr, err := utils.PrepareMapArray(m, utils.RestMethodsMAT)
	if err != nil {
		return nil, fmt.Errorf("collecting rest methods ended with an error: %w", err)
	}

	for title := range m {
		pm := m[title].([]*yaml.Node)
		if pm == nil {
			return nil, warp_errors.ErrNilMap
		}
		for i := range pm {
			switch pm[i].Value {
			case "post":
				mapArr = append(mapArr, utils.Unwrap(pm))
			case "get":
				mapArr = append(mapArr, utils.Unwrap(pm))
			}
		}

	}
	return mapArr, nil
}



func (p *Parser) NewDefinition() *models.Definitions {
	m, _ := p.collectDefinitions()
	var str []string
	for i := range m {
		for n := range m[i] {
			str = append(str, n)
		}
	}
	fmt.Println(str)
	kek := models.NewDefinition(str[0])
	return kek
}
func (p *Parser) collectDefinitions() ([]map[string]any, error) {
	m := p.SwaggerCfg.SwagBP.DefinitionsMap
	ooo := p.SwaggerCfg.Definitions.DefinitionsNode

	mapArr, err := utils.PrepareMapArray(m, utils.DefinitionsMAT)
	if err != nil {
		return nil, fmt.Errorf("collecting definitions ended with an error: %w", err)
	}
	mapArr = append(mapArr, utils.Unwrap(p.SwaggerCfg.Definitions.DefinitionsNode))
	return mapArr, nil
}
*/
