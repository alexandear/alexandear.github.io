#! /bin/sh

go install golang.org/dl/go1.21.13@latest
go1.21.13 download

go1.21.13 build ./... && go1.21.13 test -run ^$ ./...
