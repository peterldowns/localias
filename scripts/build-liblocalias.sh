#!/usr/bin/env bash
export CGO_ENABLED=1
export CC=/usr/bin/clang
export CXX=/usr/bin/clang++
rm -rf ./build && mkdir -p ./build
# amd
GOOS=darwin GOARCH=amd64 go build --buildmode=c-archive -o ./build/liblocalias-amd64.a ./app/
# arm
GOOS=darwin GOARCH=arm64 go build --buildmode=c-archive -o ./build/liblocalias-arm64.a ./app/
# smash them together
lipo -create ./build/*.a -o ./app/Localias/liblocalias.a
mv ./build/liblocalias-arm64.h ./app/Localias/liblocalias.h
rm -rf build
