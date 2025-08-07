#! /bin/sh

# https://alexandear.github.io/posts/2024-07-12-old-go-darwin-arm64/
export GOARCH=amd64
go run golang.org/dl/go1.15.15@latest download
go install golang.org/dl/go1.15.15@latest

# Capture build output and handle "no packages to build" case
build_output=$(go1.15.15 build ./... 2>&1)
build_exit_code=$?

if [ $build_exit_code -ne 0 ] && echo "$build_output" | grep -q "no packages to build"; then
    echo "No packages to build, skipping..."
    build_exit_code=0
elif [ $build_exit_code -ne 0 ]; then
    echo "$build_output"
    exit $build_exit_code
fi

go1.15.15 test -run ^$ ./...
