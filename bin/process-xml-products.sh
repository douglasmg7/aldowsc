#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Go to source path.
cd $(dirname $0)
cd ..

# Last downloaded XML file.
F_LAST=$ZUNKAPATH/xml/aldo/aldo-products-substitution.xml

if [[ $1 == dev ]]; then
    go run *.go dev < $F_LAST
else
    aldowsc < $F_LAST
fi

# read -p "Run this command only in dev mode, 'y' to continue. " answer
# # Just run on dev mode.
# [[ $answer == 'y' ]] && go run *.go dev < $F_LAST
