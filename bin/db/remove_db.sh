#!/usr/bin/env bash 

# ZUNKAPATH not defined.
if [ -z "$ZUNKAPATH" ]; then
	printf "error: ZUNKAPATH not defined.\n" >&2
	exit 1 
fi

# ZUNKA_ALDOWSC_DB not defined.
if [ -z "$ZUNKA_ALDOWSC_DB" ]; then
	printf "error: ZUNKA_ALDOWSC_DB not defined.\n" >&2
	exit 1 
fi

printf "Removing db %s/%s\n" $ZUNKAPATH/db/$ZUNKA_ALDOWSC_DB
rm $ZUNKAPATH/db/$ZUNKA_ALDOWSC_DB
