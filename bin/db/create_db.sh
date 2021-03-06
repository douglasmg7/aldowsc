#!/usr/bin/env bash 

# ZUNKAPATH not defined.
if [ -z "$ZUNKAPATH" ]; then
	printf "error: ZUNKAPATH not defined.\n" >&2
	exit 1 
fi

# ZUNKAPATH not defined.
if [ -z "$ZUNKA_ALDOWSC_DB" ]; then
	printf "error: ZUNKA_ALDOWSC_DB not defined.\n" >&2
	exit 1 
fi

DB=$ZUNKAPATH/db/$ZUNKA_ALDOWSC_DB

# Create db if not exist.
if [[ ! -f $DB ]]; then
	echo Creating $DB
    mkdir -p $ZUNKAPATH/db
    sqlite3 $DB < $(dirname $0)/tables.sql
fi
