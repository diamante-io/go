#! /usr/bin/env bash

set -e

mcdev-each-change go test {{.Pkg}}
