#!/usr/bin/env bash
if [[ $(git status --porcelain) == "" ]]; then 
    echo "generated files are up to date"
    exit 0
fi
echo "detected changes in committed and generated files, which means that you"
echo "either haven't re-generated them or have generated them based"
echo "on files you haven't committed."
echo ""
echo "you can re-generate all of these files locally by running:"
echo ""
echo "    just generate"
echo "    just lint"
echo ""
echo "the differences are:"
echo ""
git status --porcelain
echo ""
git diff
echo ""
exit 1
