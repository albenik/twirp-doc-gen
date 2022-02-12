package doc

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	maxRecurseLevel = 3
	listItemsCount  = 3
)

var (
	anyValues = [listItemsCount]proto.Message{
		wrapperspb.Int64(12345),
		wrapperspb.Bool(true),
		wrapperspb.String("string"),
	}

	stringValues = [listItemsCount]string{
		"foo",
		"bar",
		"baz",
	}
)

func fillMessageFields(msg protoreflect.Message, level, iteration int) {
	if level > maxRecurseLevel {
		return
	}

	fieldDescs := msg.Descriptor().Fields()
	for i := 0; i < fieldDescs.Len(); i++ {
		fd := fieldDescs.Get(i)
		fk := fd.Kind()

		switch {
		case fd.IsList():
			setList(msg.Mutable(fd).List(), fd, level)

		case fd.IsMap():
			setMap(msg.Mutable(fd).Map(), fd, level, iteration)

		case fk == protoreflect.MessageKind || fk == protoreflect.GroupKind:
			msg.Set(fd, messageValue(fd.Message(), level, iteration))

		default:
			msg.Set(fd, scalarValue(fd.Kind(), iteration))
		}
	}
}

func setList(list protoreflect.List, fd protoreflect.FieldDescriptor, level int) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for i := 0; i < listItemsCount; i++ {
			val := list.NewElement()
			fillMessageFields(val.Message(), level+1, i)
			list.Append(val)
		}
	default:
		for i := 0; i < listItemsCount; i++ {
			list.Append(scalarValue(fd.Kind(), i))
		}
	}
}

func setMap(pmap protoreflect.Map, fd protoreflect.FieldDescriptor, level, iteration int) {
	fields := fd.Message().Fields()
	keyDesc := fields.ByNumber(1)
	valDesc := fields.ByNumber(2)

	pkey := scalarValue(keyDesc.Kind(), iteration)

	switch kind := valDesc.Kind(); kind {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		pmap.Set(pkey.MapKey(), messageValue(valDesc.Message(), level, iteration))
	default:
		pmap.Set(pkey.MapKey(), scalarValue(kind, iteration))
	}
}

func messageValue(md protoreflect.MessageDescriptor, level, iteration int) protoreflect.Value {
	switch md.FullName() {
	case googleProtobufAny:
		any, err := anypb.New(anyValues[iteration])
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOfMessage(any.ProtoReflect())

	case googleProtobufDuration:
		return protoreflect.ValueOfMessage(durationpb.New(13 * time.Second).ProtoReflect())

	default:
		val := protoreflect.ValueOfMessage(dynamicpb.NewMessage(md))
		fillMessageFields(val.Message(), level+1, 0)
		return val
	}
}

func scalarValue(kind protoreflect.Kind, iteration int) protoreflect.Value {
	switch kind {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(true)

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(1 << 30)

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(1 << 30)

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(1 << 30)

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(1 << 30)

	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(3.14159265)

	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(3.14159265)

	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte("bytes"))

	case protoreflect.StringKind:
		return protoreflect.ValueOfString(stringValues[iteration])

	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(1)
	}

	panic(fmt.Errorf("FieldDescriptor.Kind %v is not valid", kind))
}
