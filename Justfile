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
  go build -o bin/localias ./cmd/localias

build-liblocalias:
  #!/usr/bin/env bash
  set -x
  export CGO_ENABLED=1
  export CC=/usr/bin/clang
  export CXX=/usr/bin/clang++
  rm -rf ./build
  mkdir -p ./build
  # use zig as a cross-compiler because the nix-provided clang cannot do it.
  # could also use the system-provided clang at /usr/bin/clang.
  #ZIGFLAGS="-target x86_64-macos" CXX="zig c++ $ZIGFLAGS" CC="zig cc $ZIGFLAGS" GOOS=darwin GOARCH=amd64 go build --buildmode=c-archive -o ./build/liblocalias-amd64.a ./app/
  CC=/usr/bin/clang CXX=/usr/bin/clang++ GOOS=darwin GOARCH=amd64 go build --buildmode=c-archive -o ./build/liblocalias-amd64.a ./app/
  #ZIGFLAGS="-target aarch64-macos" CXX="zig c++ $ZIGFLAGS" CC="zig cc $ZIGFLAGS" GOOS=darwin GOARCH=arm64 go build --buildmode=c-archive -o ./build/liblocalias-arm64.a ./app/
  #ZIGFLAGS="-target aarch64-macos" CXX="zig c++ $ZIGFLAGS" CC="zig cc $ZIGFLAGS"
  CC=/usr/bin/clang CXX=/usr/bin/clang++ GOOS=darwin GOARCH=arm64 go build --buildmode=c-archive -o ./build/liblocalias-arm64.a ./app/
  lipo -create ./build/*.a -o ./Localias/liblocalias.a
  mv ./build/liblocalias-arm64.h ./Localias/liblocalias.h
