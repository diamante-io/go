#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $DIR

openssl genrsa -des3 -passout pass:x -out new.pass.key 2048
openssl rsa -passin pass:x -in new.pass.key -out new.key
rm new.pass.key
openssl req -new -key new.key -out new.csr -config localhost.conf
openssl x509 -req -days 365 -in new.csr -signkey new.key -out new.crt

mv new.csr server.csr
mv new.crt server.crt
mv new.key server.key
