package proto_parser //nolint:all

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yoheimuta/go-protoparser/v4"
	"os"
)

type IProtoParser interface {
	Parse(path string, selectedMsg string) ([]*ParsedMsg, error)
}

func NewIProtoParser() IProtoParser {
	return &parser{MapProto: make(map[string][]*MessageBody)}
}

type parser struct {
	ParsedCustom []*ParsedMsg
	ParsedSimple []*ParsedMsg
	MapProto     map[string][]*MessageBody
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

func newParsedMsg(name, tYpe string, custom bool) *ParsedMsg {
	return &ParsedMsg{
		IsCustom:    custom,
		ParsedNames: name,
		ParsedTypes: tYpe,
	}
}
func (p *parser) Parse(path string, selectedMsg string) ([]*ParsedMsg, error) {
	if err := p.parse(path); err != nil {
		return nil, fmt.Errorf("failed while parsing, %w", err)
	}
	var parsedMessage []*ParsedMsg
	if len(p.MapProto) == 0 {
		return nil, errors.New("p.MapProto is empty")
	}
	for msgName, msgBody := range p.MapProto {
		if msgName == selectedMsg {
			for i := range msgBody {
				if len(msgBody) == 0 {
					return nil, errors.New("msgBody is empty")
				}
				ppType := ParseProtoType(msgBody[i].Type)
				switch ppType {
				case CustomProtoType:
					parsedMessage = append(parsedMessage, newParsedMsg(parsedMessage[i].ParsedNames, parsedMessage[i].ParsedTypes, true))
				default:

					parsedMessage = append(parsedMessage, newParsedMsg(msgBody[i].FieldName, msgBody[i].Type, false))
				}
			}
		}
	}

	return parsedMessage, nil
}
func (p *parser) parse(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open proto file: %w", err)
	}
	defer file.Close()
	got, err := protoparser.Parse(
		file,
		protoparser.WithFilename(file.Name()),
	)
	if err != nil {
		return fmt.Errorf("failed to parse %w", err)
	}
	var model *Body
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %w", err)
	}
	err = json.Unmarshal(gotJSON, &model)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %w", err)
	}
	for i := range model.ProtoBody {
		p.MapProto[model.ProtoBody[i].MessageName] = model.ProtoBody[i].MessageBody
	}
	return nil
}
