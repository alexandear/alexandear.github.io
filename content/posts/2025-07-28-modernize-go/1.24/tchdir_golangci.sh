#! /bin/sh

# << snippet begin >>
# No auto fix: os.Chdir() could be replaced by t.Chdir() in TestSomeFunc (usetesting)
golangci-lint run --no-config --enable-only usetesting --issues-exit-code 0 ./...
# << snippet end >>
