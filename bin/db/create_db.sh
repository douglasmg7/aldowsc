#!/usr/bin/env bash 

# ZUNKAPATH not defined.
if [ -z "$ZUNKAPATH" ]; then
	printf "error: ZUNKAPATH not defined.\n" >&2
	exit 1 
fi

DB_NAME=$($(dirname $0)/../get-toml-value.sh aldowsc dbFileName $ZUNKAPATH/config.toml)
DB_PATH=$ZUNKAPATH/db
DB=$DB_PATH/$DB_NAME

# Create db if not exist.
if [[ ! -f $DB ]]; then
	echo Creating $DB
    mkdir -p $DB_PATH
    sqlite3 $DB < $(dirname $0)/tables.sql
fi
