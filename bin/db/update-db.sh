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

# Updating db.
echo Updating $DB
sqlite3 $DB < $(dirname $0)/tables.sql
sqlite3 $DB < $(dirname $0)/copy-old-to-new.sql


