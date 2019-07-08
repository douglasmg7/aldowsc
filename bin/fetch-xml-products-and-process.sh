#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

[[ $ZUNKAENV == PRODUCTION ]] && printf "Can not run in production mode!\n" $$ exit 1 || go run *.go < $ZUNKAPATH/xml/test.xml


if [[ "$ZUNKAENV" == "PRODUCTION" ]]
then
	curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
		-s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
		2>>$ZUNKAPATH/log/xml_download_time.log \
		| tee $ZUNKAPATH/xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
		| aldowsc
else
	# Go to source path.
	cd $(dirname $0)
	cd ..
	curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
		-s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
		2>>$ZUNKAPATH/log/xml_download_time.log \
		| tee $ZUNKAPATH/xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
		| go run *.go
fi

# curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
	# -s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
	# 2>>log/xml_download_time.log \
	# | tee xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
	# | go run *.go