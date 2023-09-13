package proto_parser //nolint:revive

import (
	"encoding/json"
	"fmt"
	"github.com/yoheimuta/go-protoparser/v4"
	"os"
)

type IProtoParser interface {
	Parse(path string) ([]*ParsedMsg, []*ParsedMsg, error)
}

func NewIProtoParser() IProtoParser {
	return &parser{}
}

type parser struct {
	ParsedCustom []*ParsedMsg
	ParsedSimple []*ParsedMsg
}
type ParsedMsg struct {
	IsCustom                 bool
	ParsedNames, ParsedTypes string
}
type Body struct {
	ProtoBody []*ProtoBody `json:"ProtoBody"`
}
type ProtoBody struct {
	MessageName string         `json:"MessageName"`
	MessageBody []*MessageBody `json:"MessageBody"`
}
type MessageBody struct {
	IsRepeated bool   `json:"IsRepeated"`
	IsRequired bool   `json:"IsRequired"`
	IsOptional bool   `json:"IsOptional"`
	FieldName  string `json:"FieldName"`
	Type       string `json:"Type"`
}

func newParsedMsg(t, n string, custom bool) *ParsedMsg {
	return &ParsedMsg{
		IsCustom:    custom,
		ParsedNames: n,
		ParsedTypes: t,
	}
}
func (p *parser) Parse(path string) ([]*ParsedMsg, []*ParsedMsg, error) {
	parsedNames, parsedTypes, err := p.parse(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed while parsing, %w", err)
	}
	if len(parsedTypes) != len(parsedNames) {
		return nil, nil, ErrLenNotEqual
	}
	var simpleTypes []*ParsedMsg
	var customTypes []*ParsedMsg
	for i := range parsedTypes {
		ppType := ParseProtoType(parsedTypes[i])
		switch ppType {
		case CustomProtoType:
			customTypes = append(customTypes, newParsedMsg(parsedTypes[i], parsedNames[i], true))
		default:
			simpleTypes = append(simpleTypes, newParsedMsg(parsedTypes[i], parsedNames[i], true))
		}
	}
	return simpleTypes, customTypes, nil
}
func (p *parser) parse(path string) ([]string, []string, error) {
	var parsedNames, parsedTypes []string
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open proto file: %w", err)
	}
	defer file.Close()
	got, err := protoparser.Parse(
		file,
		protoparser.WithFilename(file.Name()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse %w", err)
	}

	var model *Body
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal %w", err)
	}
	err = json.Unmarshal(gotJSON, &model)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal %w", err)
	}
	for i := range model.ProtoBody {
		for ii := range model.ProtoBody[i].MessageBody {
			// fmt.Printf("%s %s \n", model.ProtoBody[i].MessageBody[ii].FieldName, model.ProtoBody[i].MessageBody[ii].Type)
			parsedNames = append(parsedNames, model.ProtoBody[i].MessageBody[ii].FieldName)
			parsedTypes = append(parsedTypes, model.ProtoBody[i].MessageBody[ii].Type)
		}
	}
	return parsedNames, parsedTypes, nil
}
