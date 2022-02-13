package doc

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	md "github.com/albenik-go/twirp-doc-gen/internal/markdown"
)

type Generator struct {
	baseURL  string
	writer   io.Writer
	doc      *md.Document
	messages map[string]*protogen.Message
	enums    map[string]*protogen.Enum
}

func NewGenerator(w io.Writer, baseURL string) *Generator {
	return &Generator{
		baseURL: baseURL,
		writer:  w,
	}
}

func (g *Generator) GenerateServiceDocument(service *protogen.Service) error {
	g.messages = make(map[string]*protogen.Message)
	g.enums = make(map[string]*protogen.Enum)

	g.doc = new(md.Document)
	g.doc.Append(md.TH1(string(service.Desc.Name())))
	g.doc.Append(md.P(md.Code(string(service.Desc.FullName()))))

	if desc := descriptionBlock(service.Comments.Leading); desc != nil {
		g.doc.Append(desc)
	}

	methodListIems := make([]md.Block, 0, len(service.Methods))
	for _, method := range service.Methods {
		g.collectModels(method.Input)
		g.collectModels(method.Output)

		mname := string(method.Desc.Name())
		methodListIems = append(methodListIems, md.LinkToHeader(fmt.Sprintf("POST /%s", mname), mname))
	}

	g.doc.Append(md.TH3("Methods"))
	g.doc.Append(md.UL(methodListIems...))

	modelsCount := len(g.messages) + len(g.enums)
	modelKeys := make([]string, 0, modelsCount)
	if modelsCount > 0 {
		for k := range g.enums {
			modelKeys = append(modelKeys, k)
		}
		for k := range g.messages {
			modelKeys = append(modelKeys, k)
		}
		sort.Strings(modelKeys)
		modelListItems := make([]md.Block, 0, len(modelKeys))
		for _, k := range modelKeys {
			var name string
			if m, ok := g.messages[k]; ok {
				name = string(m.Desc.FullName())
			} else {
				name = string(g.enums[k].Desc.FullName())
			}
			modelListItems = append(modelListItems, md.LinkToHeader(name, name))
		}
		g.doc.Append(md.TH3("Models"))
		g.doc.Append(md.UL(modelListItems...))
	}

	g.doc.Append(md.Line())

	g.doc.Append(md.TH2("Methods"))
	g.doc.Append(md.P(
		md.T("Base URL: "),
		md.Code(g.baseURL+"/"+string(service.Desc.FullName())),
	))

	for _, method := range service.Methods {
		g.doc.Append(md.H3(md.G(
			md.T("POST "),
			md.Code(fmt.Sprintf("/%s", method.Desc.Name())),
		)))
		if desc := descriptionBlock(method.Comments.Leading); desc != nil {
			g.doc.Append(desc)
		}

		g.doc.Append(md.TH4("Request"))
		g.doc.Append(md.P(md.Code(string(method.Input.Desc.FullName()))))
		g.doc.Append(md.P(md.Code(fmt.Sprintf("POST /%s/%s", service.Desc.FullName(), method.Desc.Name()))))
		g.doc.Append(md.CodeBlock(messageJSONString(method.Input.Desc), "json"))
		g.printMessageFields(method.Input)

		g.doc.Append(md.TH4("Response"))
		g.doc.Append(md.P(md.Code(string(method.Output.Desc.FullName()))))
		g.doc.Append(md.P(md.Code("HTTP 200 OK")))
		g.doc.Append(md.CodeBlock(messageJSONString(method.Output.Desc), "json"))
		g.printMessageFields(method.Output)
	}

	if modelsCount > 0 {
		g.doc.Append(md.TH2("Models"))

		for _, k := range modelKeys {
			if m, ok := g.messages[k]; ok {
				g.doc.Append(md.TH3(string(m.Desc.FullName())))
				g.printMessageFields(m)
				continue
			}

			e := g.enums[k]
			g.doc.Append(md.TH3(string(e.Desc.FullName())))
			g.printEnumItems(e)
		}
	}

	g.doc.Append(md.TH2("Twirp Errors"))
	g.doc.Append(twirpErrorCodesTable())

	return g.doc.Generate(g.writer)
}

func (g *Generator) collectModels(message *protogen.Message) {
	for _, field := range message.Fields {
		switch field.Desc.Kind() {
		case protoreflect.MessageKind:
			if !field.Desc.IsMap() {
				if _, ok := protoKnownTypeLabels[field.Message.Desc.FullName()]; ok {
					break
				}
				g.messages[string(field.Message.Desc.FullName())] = field.Message
			}
			g.collectModels(field.Message)
		case protoreflect.EnumKind:
			g.enums[string(field.Enum.Desc.FullName())] = field.Enum
		}
	}
}

func (g *Generator) printMessageFields(message *protogen.Message) {
	if desc := descriptionBlock(message.Comments.Leading); desc != nil {
		g.doc.Append(desc)
	}

	t := new(md.Table)
	t.AddColumn("Field", md.AlignLeft)
	t.AddColumn("Type", md.AlignCenter)
	t.AddColumn("Description", md.AlignLeft)

	for _, field := range message.Fields {
		fieldComment := descriptionCellText(field.Comments.Leading)
		if fieldComment == nil {
			fieldComment = md.T("")
		}
		t.AppendRow(
			md.Code(field.Desc.JSONName()),
			fieldTypeBlock(field),
			fieldComment,
		)
	}

	g.doc.Append(t)
}

func (g *Generator) printEnumItems(enum *protogen.Enum) {
	if desc := descriptionBlock(enum.Comments.Leading); desc != nil {
		g.doc.Append(desc)
	}

	table := new(md.Table)
	table.AddColumn("Value", md.AlignLeft)
	table.AddColumn("Description", md.AlignLeft)

	for _, ev := range enum.Values {
		desc1 := descriptionCellText(ev.Comments.Leading)
		desc2 := descriptionCellText(ev.Comments.Trailing)

		desc := md.T("")
		if desc1 != nil {
			if desc2 != nil {
				desc = md.G(md.CellP(desc1), md.CellP(desc2))
			} else {
				desc = desc1
			}
		} else if desc2 != nil {
			desc = desc2
		}

		table.AppendRow(md.Code(string(ev.Desc.Name())), desc)
	}

	g.doc.Append(table)
}

func messageJSONString(mdesc protoreflect.MessageDescriptor) string {
	m := dynamicpb.NewMessage(mdesc)
	fillMessageFields(m, 0, 0)

	j := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	return j.Format(m)
}

func fieldTypeBlock(field *protogen.Field) md.Block {
	var block md.Block

	//nolint:exhaustive
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		msg := field.Message
		name := msg.Desc.FullName()

		if s, ok := protoKnownTypeLabels[name]; ok {
			block = md.T(s)
			break
		}

		if field.Desc.IsMap() {
			var key, val *protogen.Field
			for _, f := range field.Message.Fields {
				if f.Desc == field.Desc.MapKey() {
					key = f
				} else {
					val = f
				}
			}
			block = md.G(
				md.T("map "),
				fieldTypeBlock(key),
				md.T(" to "),
				fieldTypeBlock(val),
			)
			break
		}

		block = md.LinkToHeader(string(name), string(name))

	case protoreflect.EnumKind:
		enum := field.Desc.Enum()
		block = md.Link(string(enum.Name()), string(enum.FullName()))

	default:
		if s, ok := protoKindTypes[field.Desc.Kind()]; ok {
			block = md.T(s)
		}
	}

	if block == nil {
		panic(fmt.Errorf("field %s: unknown kind %s", field.Desc.FullName(), field.Desc.Kind()))
	}

	if field.Desc.IsList() {
		return md.G(
			md.T("array of "),
			block,
		)
	}
	return block
}

func descriptionBlock(c protogen.Comments) md.Block {
	if c == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(c)), "\n\n")

	blocks := make([]md.Block, 0, len(lines))
	for _, l := range lines {
		blocks = append(blocks, md.P(md.TI(strings.ReplaceAll(strings.TrimSpace(l), "\n ", "\n"))))
	}

	return md.G(blocks...)
}

func descriptionCellText(c protogen.Comments) md.Block {
	s := string(c)
	if s == "" {
		return nil
	}
	return md.T(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\n\n", "<br/>"), "\n", " ")))
}

func twirpErrorCodesTable() md.Block {
	t := new(md.Table)

	t.AddColumn("Twirp Error Code", md.AlignLeft)
	t.AddColumn("HTTP Status", md.AlignCenter)
	t.AddColumn("Description", md.AlignLeft)

	t.AppendRow(md.Code("invalid_argument"), md.Code("400"), md.T("The client specified an invalid argument. "+
		"This indicates arguments that are invalid regardless of the state of the system "+
		"(i.e. a malformed file name, required argument, number out of range, etc.)."))
	t.AppendRow(md.Code("malformed"), md.Code("400"), md.T("The client sent a message which could not be decoded. "+
		"This may mean that the message was encoded improperly or that the client and server have incompatible "+
		"message definitions."))
	t.AppendRow(md.Code("out_of_range"), md.Code("400"), md.T("The operation was attempted past the valid range. "+
		"For example, seeking or reading past end of a paginated collection. Unlike \"invalid_argument\", "+
		"this error indicates a problem that may be fixed if the system state changes "+
		"(i.e. adding more items to the collection). There is a fair bit of overlap between "+
		"\"failed_precondition\" and \"out_of_range\". We recommend using \"out_of_range\" "+
		"(the more specific error) when it applies so that callers who are iterating through a space"+
		" can easily look for an \"out_of_range\" error to detect when they are done."))
	t.AppendRow(md.Code("unauthenticated"), md.Code("401"), md.T("The request does not have valid authentication "+
		"credentials for the operation."))
	t.AppendRow(md.Code("permission_denied"), md.Code("403"), md.T("The caller does not have permission to execute "+
		"the specified operation. It must not be used if the caller cannot be identified "+
		"(use \"unauthenticated\" instead)."))
	t.AppendRow(md.Code("bad_route"), md.Code("404"), md.T("The requested URL path wasn't routable to a Twirp "+
		"service and method. This is returned by generated server code and should not be returned by "+
		"application code (use \"not_found\" or \"unimplemented\" instead)."))
	t.AppendRow(md.Code("not_found"), md.Code("404"), md.T("Some requested entity was not found."))
	t.AppendRow(md.Code("canceled"), md.Code("408"), md.T("The operation was cancelled."))
	t.AppendRow(md.Code("deadline_exceeded"), md.Code("408"), md.T("Operation expired before completion. "+
		"For operations that change the state of the system, this error may be returned even if the operation "+
		"has completed successfully (timeout)."))
	t.AppendRow(md.Code("already_exists"), md.Code("409"), md.T("An attempt to create an entity failed because one "+
		"already exists."))
	t.AppendRow(md.Code("aborted"), md.Code("409"), md.T("The operation was aborted, typically due to "+
		"a concurrency issue like sequencer check failures, transaction aborts, etc."))
	t.AppendRow(md.Code("failed_precondition"), md.Code("412"), md.T("The operation was rejected because the "+
		"system is not in a state required for the operation's execution. For example, doing an rmdir "+
		"operation on a directory that is non-empty, or on a non-directory object, or when having conflicting "+
		"read-modify-write on the same resource."))
	t.AppendRow(md.Code("resource_exhausted"), md.Code("429"), md.T("Some resource has been exhausted or "+
		"rate-limited, perhaps a per-user quota, or perhaps the entire file system is out of space."))
	t.AppendRow(md.Code("unknown"), md.Code("500"), md.T("An unknown error occurred. For example, "+
		"this can be used when handling errors raised by APIs that do not return any error information."))
	t.AppendRow(md.Code("internal"), md.Code("500"), md.T("When some invariants expected by the underlying system "+
		"have been broken. In other words, something bad happened in the library or backend service. "+
		"Twirp specific issues like wire and serialization problems are also reported as \"internal\" errors."))
	t.AppendRow(md.Code("unavailable"), md.Code("500"), md.T("The service is currently unavailable. This is most "+
		"likely a transient condition and may be corrected by retrying with a backoff."))
	t.AppendRow(md.Code("unimplemented"), md.Code("501"), md.T("The operation is not implemented or not "+
		"supported/enabled in this service."))
	t.AppendRow(md.Code("dataloss"), md.Code("503"), md.T("The operation resulted in unrecoverable data loss or "+
		"corruption."))

	return md.G(
		md.P(md.Link("https://twitchtv.github.io/twirp/docs/spec_v7.html#error-codes", "Official documentation")),
		t,
	)
}
