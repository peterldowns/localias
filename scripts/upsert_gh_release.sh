#!/usr/bin/env bash
# Attempts to get an upload URL for adding assets to the current
# Github Release. It does this by attemping to create the release,
# then reading it if it already exists. 
#
# On success, the result is a 0 status code and a valid upload URL.
#
# Docs:
#   https://docs.github.com/en/rest/releases/releases#create-a-release
#   https://stackoverflow.com/questions/45240336/how-to-use-github-release-api-to-make-a-release-without-source-code
version=$(cat ./VERSION)
commit=$(git rev-parse --short HEAD || echo 'unknown')
# https://semver.org/#spec-item-10
release_name="$version+commit.$commit"
echo "attempting to create release $release_name"
upload_url=$(
  gh api --method POST 'repos/{owner}/{repo}/releases' \
    -F "tag_name=testing-$release_name" \
    -F "name=testing-$release_name" \
    -F "target_comitish=$commit_sha" \
    --jq '.upload_url' \
)
if [ $? != 0 ]; then 
  echo "failed to create release:"
  echo "$upload_url"
  echo "attempting to fetch existing release $release_name"
  upload_url=$(
    gh api --method GET 'repos/{owner}/{repo}/releases/tags/'"$release_name" \
      --jq '.upload_url' \
  )
  if [ $? != 0 ]; then 
    echo "failed to fetch existing release $release_name"
    echo "upload_url"
    exit 1
  fi
fi
# the upload url looks like
#   https://uploads.github.com/.../<release_id>/assets{?name,label}
# this trick strips off the {?name,label}
upload_url="${upload_url%\{*}"
echo $upload_url
