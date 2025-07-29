#! /bin/sh

go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -category forvar --fix forloop_before.go
