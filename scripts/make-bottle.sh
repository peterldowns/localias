#!/usr/bin/env bash
export VERSION=$(cat ./VERSION)
export COMMIT="$(git rev-parse --short HEAD || echo 'unknown')"
export TAG_NAME="$VERSION+commit.$COMMIT"
export MODIFIED_TIME=$(date +%s)
export BINARY=$1 #./bin/localias-darwin-arm64
export BOTTLE_TAG=$2 # arm64_monterey
export TAP_GIT_HEAD=$(gh api --method GET 'repos/peterldowns/homebrew-tap/commits/HEAD' | jq -r .sha)
export LDFLAGS=$(./scripts/golang-ldflags.sh $VERSION $COMMIT)
export SOURCE_URL="https://github.com/peterldowns/localias/archive/refs/tags/$VERSION+commit.$COMMIT.tar.gz"
wget "$SOURCE_URL" -O release.tar.gz
export SOURCE_SHA_256=$(shasum -a 256 release.tar.gz | cut -d ' ' -f 1)
export BOTTLE_ROOT="https://github.com/peterldowns/localias/releases/download/$VERSION%2Bcommit.$COMMIT"


_BOTTLEDIR="localias/$VERSION/"
_BREWDIR="localias/$VERSION/.brew"
_BINDIR="localias/$VERSION/bin"
mkdir -p "$_BOTTLEDIR"
mkdir -p "$_BREWDIR"
mkdir -p "$_BINDIR"

# copy the binary
cp "$BINARY" "$_BINDIR/localias"
# copy the metadata files
cp ./README.md "$_BOTTLEDIR/README.md"
cp ./LICENSE "$_BOTTLEDIR/LICENSE"
# render the templates
envsubst < ./.brew/INSTALL_RECEIPT.tpl.json > "$_BOTTLEDIR/INSTALL_RECEIPT.json"
envsubst < ./.brew/localias.tpl.rb > "$_BREWDIR/localias.rb"

export BOTTLE_NAME="localias-$VERSION.${BOTTLE_TAG}.bottle.tar.gz"
tar -czf "$BOTTLE_NAME" ./localias
echo $BOTTLE_NAME
