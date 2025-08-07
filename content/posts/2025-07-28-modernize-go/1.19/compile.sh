#! /bin/sh

go install golang.org/dl/go1.19.13@latest
go1.19.13 download

go1.19.13 build ./... && go1.19.13 test -run ^$ ./...
