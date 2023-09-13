package proto_parser // nolint:all

import (
	"errors"
	"strings"
)

type ProtoType int

const (
	ProtoDouble ProtoType = iota
	ProtoFloat
	ProtoInt32
	ProtoInt64
	ProtoUint32
	ProtoUint64
	ProtoSint32
	ProtoSint64
	ProtoFixed32
	ProtoFixed64
	ProtoSfixed32
	ProtoSfixed64
	ProtoBool
	ProtoString
	ProtoBytes

	// CustomProtoType
	/*
		message CustomMsg{
			string someString = 1;
		}

		message Parented{
			CustomMsg parentedMsg = 1;
		}
	*/
	CustomProtoType
)

var protoTypes = [...]string{
	ProtoDouble:   "double",
	ProtoFloat:    "float",
	ProtoInt32:    "int32",
	ProtoInt64:    "int64",
	ProtoUint32:   "uint32",
	ProtoUint64:   "uint64",
	ProtoSint32:   "sint32",
	ProtoSint64:   "sint64",
	ProtoFixed32:  "fixed32",
	ProtoFixed64:  "fixed64",
	ProtoSfixed32: "sfixed32",
	ProtoSfixed64: "sfixed64",
	ProtoBool:     "bool",
	ProtoString:   "string",
	ProtoBytes:    "bytes",
}

func (s ProtoType) String() string {
	return protoTypes[s]
}

func ParseProtoType(s string) ProtoType {
	for i, r := range protoTypes {
		if strings.ToLower(s) == r {
			return ProtoType(i)
		}
	}
	return CustomProtoType
}

var ErrLenNotEqual = errors.New("types length not equal to names length")
