#!/bin/bash -e

# DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# cd $DIR

# if [ -d "$DIR/docs" ]; then
#     rm -rf $DIR/docs/*
# else
#     mkdir -p $DIR/docs
# fi

# FILES=(wiki/*)
# for f in wiki/*.md
# do
#     grip --gfm --context=patchwork-toolkit/patchwork --export $f "docs/$(basename $f .md).html"
# done