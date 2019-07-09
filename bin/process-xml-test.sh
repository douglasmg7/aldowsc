#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Go to source path.
cd $(dirname $0)
cd ..

# Just run on dev mode.
[[ $ZUNKAENV == PRODUCTION ]] && printf "Can not run in production mode!\n" $$ exit 1 || go run *.go < $ZUNKAPATH/xml/test.xml