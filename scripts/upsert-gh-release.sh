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
# Read base version from ./VERSION (e.g. v3.0.1) and count commits since ./VERSION
# was last modified to increment the patch between VERSION edits.
base_version=$(tr -d '[:space:]' < ./VERSION)
version_commit=$(git log -n 1 --pretty=format:%H -- ./VERSION 2>/dev/null || true)
if [ -n "$version_commit" ]; then
  commits_since=$(git rev-list --count "${version_commit}..HEAD")
else
  commits_since=0
fi
# Split "v3.0.1" into "v3.0" (%.*) and "1" (##*.), then offset patch by commit count
major_minor="${base_version%.*}"
base_patch="${base_version##*.}"
patch=$((base_patch + commits_since))

commit_sha=$(git rev-parse --verify HEAD) || exit 1
release_name="${major_minor}.${patch}"
if ! upload_url=$(
  gh api --method POST 'repos/{owner}/{repo}/releases' \
    -F "tag_name=$release_name" \
    -F "name=$release_name" \
    -F "target_commitish=$commit_sha" \
    --jq '.upload_url' \
); then 
  if ! upload_url=$(
    gh api --method GET 'repos/{owner}/{repo}/releases/tags/'"$release_name" \
      --jq '.upload_url' \
  ); then 
    exit 1
  fi
fi

# Creating a release can also succeed for an existing tag. Resolve the remote
# tag (including annotated tags) before returning any upload URL, so a VERSION
# edit or a concurrent release cannot cause uploads to another commit's release.
if ! tag_commit=$(
  gh api --method GET 'repos/{owner}/{repo}/commits/refs/tags/'"$release_name" \
    --jq '.sha' \
); then
  exit 1
fi

# The release tag must point to the commit we're building; otherwise we'd
# upload this commit's binaries to an older commit's release.
if [ "$tag_commit" != "$commit_sha" ]; then
  echo "Release tag $release_name points to $tag_commit, not $commit_sha. Increase VERSION to an unused version." >&2
  exit 1
fi

# the upload url looks like
#   https://uploads.github.com/.../<release_id>/assets{?name,label}
# this trick strips off the {?name,label}
upload_url="${upload_url%\{*}"
echo "$upload_url"
