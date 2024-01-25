#! /bin/bash
set -e

go mod tidy
git diff --exit-code -- go.mod || (echo "Go file go.mod is dirty, update the file with 'go mod tidy' locally." && exit 1)
git diff --exit-code -- go.sum || (echo "Go file go.sum is dirty, update the file with 'go mod tidy' locally." && exit 1)
diff -u go.list <(./golist.sh) || (echo "Go dependencies have changed, update the go.list file with './golist.sh > go.list' locally." && exit 1)
go mod verify || (echo "One or more Go dependencies failed verification. Either a version is no longer available, or the author or someone else has modified the version so it no longer points to the same code." && exit 1)
