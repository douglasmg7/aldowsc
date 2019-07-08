#!/usr/bin/env bash
DB_NAME="aldowsc.db" 
if [ -z "$ZUNKAPATH" ]
then
	printf "error: ZUNKAPATH enviorment not defined.\n" >&2
	exit 1 
else
	printf "Removing db %s/%s\n" $ZUNKAPATH $DB_NAME
fi
rm $ZUNKAPATH/db/$DB_NAME
