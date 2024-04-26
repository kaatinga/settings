#!/bin/bash

echo "updating direct dependencies..."
go list -m -u -f '{{if not (or .Indirect .Main)}}{{.Path}}@latest{{end}}' all | xargs go get
echo "tidying up..."
go mod tidy -compat=1.18