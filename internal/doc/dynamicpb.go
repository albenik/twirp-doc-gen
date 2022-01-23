package doc

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

const maxRecurseLevel = 3

func fillMessage(m protoreflect.Message, level int) {
	if level > maxRecurseLevel {
		return
	}

	fieldDescs := m.Descriptor().Fields()
	for i := 0; i < fieldDescs.Len(); i++ {
		fd := fieldDescs.Get(i)
		switch {
		case fd.IsList():
			setList(m.Mutable(fd).List(), fd, level)
		case fd.IsMap():
			setMap(m.Mutable(fd).Map(), fd, level)
		default:
			setScalarField(m, fd, level)
		}
	}
}

func setList(list protoreflect.List, fd protoreflect.FieldDescriptor, level int) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for i := 0; i < 3; i++ {
			val := list.NewElement()
			fillMessage(val.Message(), level+1)
			list.Append(val)
		}
	default:
		for i := 0; i < 3; i++ {
			list.Append(scalarField(fd.Kind()))
		}
	}
}

func setMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor, level int) {
	fields := fd.Message().Fields()
	keyDesc := fields.ByNumber(1)
	valDesc := fields.ByNumber(2)

	pkey := scalarField(keyDesc.Kind())
	switch kind := valDesc.Kind(); kind {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		val := mmap.NewValue()
		fillMessage(val.Message(), level+1)
		mmap.Set(pkey.MapKey(), val)
	default:
		mmap.Set(pkey.MapKey(), scalarField(kind))
	}
}

func setScalarField(m protoreflect.Message, fd protoreflect.FieldDescriptor, level int) {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		val := m.NewField(fd)
		fillMessage(val.Message(), level+1)
		m.Set(fd, val)
	default:
		m.Set(fd, scalarField(fd.Kind()))
	}
}

func scalarField(kind protoreflect.Kind) protoreflect.Value {
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
		return protoreflect.ValueOfString("string")

	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(1)
	}

	panic(fmt.Errorf("FieldDescriptor.Kind %v is not valid", kind))
}
