#!/usr/bin/env bash
# Generates the linker flags needed to embed version information in the built
# CLI binaries.
#
# The result of this script is like
#
#   -X main.Version=0.0.6 -X main.Commit=19040ae
#
# and should be used like this:
#
#   ldflags=$(./scripts/golang-ldflags.sh)
#   go build -ldflags "$ldflags" ...
#
VERSION=$1
COMMIT=$2
if [ -z "$VERSION" ]; then
  VERSION=$(cat ./VERSION)
fi
if [ -z "$COMMIT" ]; then
  COMMIT="$(git rev-parse --short HEAD || echo 'unknown')"
fi

package="github.com/peterldowns/localias/cmd/localias/shared"
echo "-X $package.Version=$VERSION -X $package.Commit=$COMMIT"
