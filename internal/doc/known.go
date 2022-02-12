package doc

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	googleProtobufAny         = "google.protobuf.Any"
	googleProtobufStringValue = "google.protobuf.StringValue"
	googleProtobufBytesValue  = "google.protobuf.BytesValue"
	googleProtobufBoolValue   = "google.protobuf.BoolValue"
	googleProtobufInt32Value  = "google.protobuf.Int32Value"
	googleProtobufInt64Value  = "google.protobuf.Int64Value"
	googleProtobufUInt32Value = "google.protobuf.UInt32Value"
	googleProtobufUInt64Value = "google.protobuf.UInt64Value"
	googleProtobufFloatValue  = "google.protobuf.FloatValue"
	googleProtobufDoubleValue = "google.protobuf.DoubleValue"
	googleProtobufTimestamp   = "google.protobuf.Timestamp"
	googleProtobufDuration    = "google.protobuf.Duration"
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

	protoKnownTypeLabels = map[protoreflect.FullName]string{
		googleProtobufAny:         googleProtobufAny,
		googleProtobufStringValue: "nullable string",
		googleProtobufBytesValue:  "bytes as nullable base64 string",
		googleProtobufBoolValue:   "nullable bool",
		googleProtobufInt32Value:  "nullable int32",
		googleProtobufInt64Value:  "nullable int64",
		googleProtobufUInt32Value: "nullable uint32",
		googleProtobufUInt64Value: "uint64 as nullable numeric string",
		googleProtobufFloatValue:  "nullable float",
		googleProtobufDoubleValue: "nullable double",
		googleProtobufTimestamp:   "datetime as nullable RFC3339 string",
		googleProtobufDuration:    "duration as nullable string",
	}
)
