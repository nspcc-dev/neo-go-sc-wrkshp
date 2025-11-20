#!/bin/sh

if [ -z "$1" ]; then
	echo "Usage: ./update_deps.sh <revision>"
	exit 1
fi

REV="$1"
root="$(git rev-parse --show-toplevel)"

for dir in "$root"/*/; do
  if ! [ -f $dir"/go.mod" ]; then continue; fi
	cd "$dir" || exit 1
	go get github.com/nspcc-dev/neo-go/pkg/interop@"$REV"
	go get github.com/nspcc-dev/neo-go@"$REV"
	go mod tidy
done
