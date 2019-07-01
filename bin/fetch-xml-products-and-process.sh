#!/usr/bin/env bash
cd $(dirname $0)
cd ..
curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" \
	-s -w "%{stderr}$(date '+%Y-%h-%d %T') - Time to download xml file: %{time_total}s\n" \
	2>>log/xml_download_time.log \
	| tee xml/aldo_products_$(date +%Y-%h-%d-%H%M%S).xml \
	| go run *.go