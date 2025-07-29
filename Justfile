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

# remove any previously-built binaries
clean:
  rm -rf ./bin
  rm -rf ./result

# run the test suite
test *args='./...':
  go test "$@"

# lint all
lint:
  just lint-go
  just lint-nix
# lint/fix go code with golangci-lint
lint-go *args='./...':
  go mod tidy
  golangci-lint config verify --config .golangci.yaml
  golangci-lint run --fix --config .golangci.yaml $@
# lint/fix nix code with nixpkgs-fmt
lint-nix:
  git ls-files '*.nix' | nixpkgs-fmt

# build the localias cli
build:
  #!/usr/bin/env bash
  ldflags=$(./scripts/golang-ldflags.sh)
  go build -ldflags "$ldflags" -o bin/localias ./cmd/localias

