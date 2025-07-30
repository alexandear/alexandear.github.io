#! /bin/sh

# << snippet begin >>
golangci-lint run --no-config --enable-only intrange --fix ./...
# << snippet end >>
