#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Go to source path.
cd $(dirname $0)

# Last downloaded XML file.
F_LAST=$ZUNKAPATH/xml/aldo/aldo-products.xml
# File to write.
F_OUT=$ZUNKAPATH/xml/aldo/aldo-products-substitution.xml

# Substitutions.
sed -f ./substitution-list.txt $F_LAST > $F_OUT