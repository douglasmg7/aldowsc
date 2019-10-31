#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Create dir if not exist.
mkdir -p $ZUNKAPATH/xml

echo $1
if [[ $1 =~ dev ]]
then
    echo Running in development mode.
	# Go to source path.
	cd $(dirname $0)
	cd ..

	# Get xml and process.
	curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
		-s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
		2>>$ZUNKAPATH/log/xml_download_time.log \
		| tee $ZUNKAPATH/xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
		| go run *.go
else
    echo Running in production mode.
	# Get xml and process.
	curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
		-s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
		2>>$ZUNKAPATH/log/xml_download_time.log \
		| tee $ZUNKAPATH/xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
		| aldowsc
fi