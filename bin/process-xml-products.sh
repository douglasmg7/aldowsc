#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Go to source path.
cd $(dirname $0)
cd ..

# Last downloaded XML file.
FILE=$ZUNKAPATH/xml/aldo/aldo-products-substitution.xml

if [[ $RUN_MODE == production ]]; then
    RUN_MODE=production aldowsc < $FILE
else
    go build
    ./aldowsc < $FILE
fi
