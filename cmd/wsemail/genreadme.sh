#! /bin/sh

set -e
set -u

go build
export PATH=`pwd`:$PATH
echo "## usage" > README.md
echo '```' >> README.md
set +e
wsemail -help 2>> README.md
set -e
echo '```' >> README.md
