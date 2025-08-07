#! /bin/sh

go install golang.org/dl/go1.24.5@latest
go1.24.5 download

# Capture build output and handle "no packages to build" case
build_output=$(go1.24.5 build ./... 2>&1)
build_exit_code=$?

if [ $build_exit_code -ne 0 ] && echo "$build_output" | grep -q "no packages to build"; then
    echo "No packages to build, skipping..."
    build_exit_code=0
elif [ $build_exit_code -ne 0 ]; then
    echo "$build_output"
    exit $build_exit_code
fi

go1.24.5 test -run ^$ ./...
