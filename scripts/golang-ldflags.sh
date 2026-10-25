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
  base_version=$(tr -d '[:space:]' < ./VERSION)
  version_commit=$(git log -n 1 --pretty=format:%H -- ./VERSION 2>/dev/null || true)
  if [ -n "$version_commit" ]; then
    commits_since=$(git rev-list --count "${version_commit}..HEAD")
  else
    commits_since=0
  fi
  major_minor="${base_version%.*}"
  base_patch="${base_version##*.}"
  patch=$((base_patch + commits_since))
  VERSION="${major_minor}.${patch}"
fi
if [ -z "$COMMIT" ]; then
  COMMIT="$(git rev-parse --short HEAD || echo '<dev>')"
fi

package="github.com/peterldowns/localias/cmd/localias/shared"
echo "-X $package.Version=$VERSION -X $package.Commit=$COMMIT"
