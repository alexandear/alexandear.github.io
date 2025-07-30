#! /bin/sh

# << snippet begin >>
golangci-lint run --no-config --enable-only copyloopvar --fix ./...
# << snippet end >>

# workaround to remove the empty line
go run mvdan.cc/gofumpt@latest -w .
