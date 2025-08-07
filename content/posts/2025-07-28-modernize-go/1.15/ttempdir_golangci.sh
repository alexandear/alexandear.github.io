#! /bin/sh

# workaround because usetesting can't detect functions from ioutil
# https://github.com/ldez/usetesting/issues/12
sed -i '' 's/ioutil\.TempDir(\([^,]*\),\s*\([^)]*\))/os.MkdirTemp(\1, \2)/g' *.go
go run golang.org/x/tools/cmd/goimports@latest -w .

# << snippet begin >>
# No auto fix: os.MkdirTemp() could be replaced by t.TempDir() in TestSomeFunc
golangci-lint run --no-config --enable-only usetesting --issues-exit-code 0 ./...
# << snippet end >>
