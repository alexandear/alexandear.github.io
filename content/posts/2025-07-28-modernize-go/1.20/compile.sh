#! /bin/sh

go install golang.org/dl/go1.20.14@latest
go1.20.14 download

go1.20.14 build ./... && go1.20.14 test -run ^$ ./...
