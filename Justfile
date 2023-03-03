# This Justfile contains rules/targets/scripts/commands that are used when
# developing. Unlike a Makefile, running `just <cmd>` will always invoke
# that command. For more information, see https://github.com/casey/just
#
#
# this setting will allow passing arguments through to tasks, see the docs here
# https://just.systems/man/en/chapter_24.html#positional-arguments
set positional-arguments

# print all available commands by default
default:
  just --list

# run the test suite
test *args='./...':
  go test "$@"

# lint the entire codebase
lint *args:
  golangci-lint run --fix --config .golangci.yaml "$@"

build:
  go build -o bin/localias .

release-binaries:
  #!/usr/bin/env bash
  GOOS=darwin GOARCH=amd64 go build -o ./bin/localias-darwin-amd64 .
  GOOS=darwin GOARCH=arm64 go build -o ./bin/localias-darwin-arm64 .
  GOOS=linux GOARCH=amd64 go build -o ./bin/localias-linux-amd64 .
  GOOS=linux GOARCH=arm64 go build -o ./bin/localias-linux-arm64 .
  commit_sha="$(git rev-parse --short HEAD)"
  timestamp="$(date +%s)"
  release_name="release-$timestamp-$commit_sha"
  token="$GITHUB_TOKEN"
  upload_url=$(curl -s -H "Authorization: token $token" \
    -X POST \
    -d "{\"tag_name\": \"$release_name\", \"name\":\"$release_name\",\"target_comitish\": \"$commit_sha\"}" \
    "https://api.github.com/repos/peterldowns/localias/releases" | jq -r '.upload_url')
  upload_url="${upload_url%\{*}"
  echo "upload_url: $upload_url"
  curl -s -H "Authorization: token $token" \
    -H "Content-Type: application/octet-stream" \
    --data-binary @bin/localias-darwin-amd64 \
    "$upload_url?name=localias-darwin-amd64&label=localias-darwin-amd64"
  curl -s -H "Authorization: token $token" \
    -H "Content-Type: application/octet-stream" \
    --data-binary @bin/localias-darwin-arm64 \
    "$upload_url?name=localias-darwin-arm64&label=localias-darwin-arm64"
  curl -s -H "Authorization: token $token" \
    -H "Content-Type: application/octet-stream" \
    --data-binary @bin/localias-linux-amd64 \
    "$upload_url?name=localias-linux-amd64&label=localias-linux-amd64"
  curl -s -H "Authorization: token $token" \
    -H "Content-Type: application/octet-stream" \
    --data-binary @bin/localias-linux-arm64 \
    "$upload_url?name=localias-linux-arm64&label=localias-linux-arm64"
