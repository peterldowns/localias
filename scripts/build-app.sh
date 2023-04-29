#!/usr/bin/env bash
pushd app || exit
rm -rf build
mkdir build
LD=clang xcodebuild -scheme Release archive -archivePath build | xcpretty
mv build.xcarchive/Products/Applications/* build
rm -rf build.xcarchive
pushd build || exit
zip -r Localias.app.zip Localias.app
popd || exit
readlink -f build/Localias.app.zip
