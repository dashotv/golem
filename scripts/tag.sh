#!/usr/bin/env bash

tag="$1"
if [[ -z $tag ]]; then
  echo "Usage: $0 <tag>"
  exit 1
fi

if [[ "$tag" == $(git tag -l "$tag") ]]; then
  echo "Tag $tag already exists"
  exit 1
fi

if [[ -n "$(git status -s)" ]]; then
  echo "There are pending changes, commit first"
  exit 1
fi

cat >config/version.go <<EOF
package config

// VERSION is the version of the application
// automatically managed by pre-commit githook
const VERSION = "$tag"
EOF
git add config/version.go
git commit -m "$tag"
git tag -a "$tag" -m "$tag"
