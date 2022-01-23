package doc

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Data types mappings
// https://swagger.io/specification/v2/
// https://developers.google.com/protocol-buffers/docs/proto3#json
var (
	protoKindTypes = map[protoreflect.Kind]string{
		protoreflect.Int32Kind:  "int32",
		protoreflect.Sint32Kind: "int32",
		protoreflect.Uint32Kind: "uint32",
		protoreflect.Int64Kind:  "int64 as numeric string",
		protoreflect.Sint64Kind: "int64 as numeric string",
		protoreflect.Uint64Kind: "uint64 as numeric string",
		protoreflect.FloatKind:  "float",
		protoreflect.DoubleKind: "double",
		protoreflect.BoolKind:   "bool",
		protoreflect.StringKind: "string",
		protoreflect.BytesKind:  "bytes as base64 string",
	}

	protoWellknownTypes = map[protoreflect.FullName]string{
		"google.protobuf.StringValue": "string",
		"google.protobuf.BytesValue":  "bytes as base64 string",
		"google.protobuf.BoolValue":   "bool",
		"google.protobuf.Int32Value":  "int32",
		"google.protobuf.Int64Value":  "int64",
		"google.protobuf.UInt32Value": "uint32",
		"google.protobuf.UInt64Value": "uint64 as numeric string",
		"google.protobuf.FloatValue":  "float",
		"google.protobuf.DoubleValue": "double",
		"google.protobuf.Timestamp":   "datetime as RFC3339 string",
		"google.protobuf.Duration":    "string",
	}
)
