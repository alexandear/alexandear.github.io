#! /bin/sh

go install golang.org/dl/go1.22.12@latest
go1.22.12 download

go1.22.12 build ./... && go1.22.12 test -run ^$ ./...
