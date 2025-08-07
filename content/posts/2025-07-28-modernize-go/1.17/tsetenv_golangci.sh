#! /bin/sh

# << snippet begin >>
# No auto fix: os.Setenv() could be replaced by t.Setenv() in TestSomeFunc
golangci-lint run --no-config --enable-only usetesting --issues-exit-code 0 ./...
# << snippet end >>
