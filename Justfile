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
  find . -name '*.nix' | xargs nixpkgs-fmt

# build the localias cli
build:
  #!/usr/bin/env bash
  VERSION=$(cat ./VERSION)
  COMMIT="$(git rev-parse --short HEAD)"
  go build -o bin/localias -ldflags \
    "-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT'" \
    ./cmd/localias

# build the localias.a library for swift app
build-liblocalias:
  ./scripts/build-liblocalias.sh

# build the swift app Localias.app
build-app:
  ./scripts/build-app.sh
