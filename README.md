# Twirp Markdown Documentation Generator

A Twirp Markdown documentation generator implemented as `protoc` plugin

Plugin generated one `<RPCServiceName>.md` file per rpc service.

# Usage

Installing the generator for protoc or buf.build:

```
go install github.com/albenik/twirp-doc-gen/cmd/protoc-gen-twirp-doc@latest
```

## Run whith the protoc

```
protoc -twirp-doc_out=path/to/doc/folder -twirp-doc_opt=paths=source_relative twirp/service/v1/service.proto
```

## Run with [buf.build](https://buf.build):

`buf.gen.yaml`:

```yaml
version: v1
plugins:
  - name: twirp-doc
    out: path/to/doc/folder
    opt: paths=source_relative
```
