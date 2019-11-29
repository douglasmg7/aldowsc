#!/usr/bin/env bash

[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 
[[ -z "$GS" ]] && printf "error: GS enviorment not defined.\n" >&2 && exit 1 

# Go to source path.
cd $(dirname $0)

# Last downloaded XML file.
F_LAST=$ZUNKAPATH/xml/aldo/aldo-products.xml
# File to write.
F_OUT=$ZUNKAPATH/xml/aldo/aldo-products-substitution.xml
# Create XML with Substitutions.
F_SUB=$GS/aldowsc/bin/substitution-list.txt

# Substitutions.
sed -f $F_SUB $F_LAST > $F_OUT