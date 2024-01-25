#!/usr/bin/env bash

go list -f '{{with .Module}}{{.Path}} {{.Version}}{{end}}' all | LC_ALL=C sort -u
