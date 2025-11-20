#!/bin/sh

if [ -z "$1" ]; then
	echo "Usage: ./update_go.sh <go-version>"
	exit 1
fi

V="$1"
root="$(git rev-parse --show-toplevel)"

for dir in "$root"/*/; do
  if ! [ -f "$dir"/go.mod ]; then continue; fi
	cd "$dir" || exit 1
	go mod edit -go="$V"
	go mod tidy
done
