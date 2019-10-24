#!/usr/bin/env bash 
DB_NAME="aldowsc.db" 
if [ -z "$ZUNKAPATH" ]
then
	printf "error: ZUNKAPATH enviorment not defined.\n" >&2
	exit 1 
else
	printf "Creating db %s/%s\n" $ZUNKAPATH/db/$DB_NAME
fi
mkdir -p $ZUNKAPATH/db
sqlite3 $ZUNKAPATH/db/$DB_NAME < tables.sql
