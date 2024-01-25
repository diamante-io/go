#!/usr/bin/env bash

set -e
set -x

source /etc/profile
# work within the current docker working dir
if [ ! -f "./diamcircle-core.cfg" ]; then
   cp /diamcircle-core.cfg ./
fi   

echo "using config:"
cat diamcircle-core.cfg

# initialize new db
diamcircle-core new-db

if [ "$1" = "standalone" ]; then
  # initialize for new history archive path, remove any pre-existing on same path from base image
  rm -rf ./history
  diamcircle-core new-hist vs

  # serve history archives to aurora on port 1570
  pushd ./history/vs/
  python3 -m http.server 1570 &
  popd
fi

exec diamcircle-core run
