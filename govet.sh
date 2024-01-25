#! /bin/bash
set -e

printf "Running go vet...\n"
go vet -all -composites=false -unreachable=false -tests=false ./...

printf "Running go vet shadow...\n"
command -v shadow >/dev/null 2>&1 || (
  dir=$(mktemp -d)
  pushd $dir
  go mod init shadow
  go get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
  popd
)

go vet -vettool=$(which shadow) ./...
