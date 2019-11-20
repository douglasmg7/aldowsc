#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# Create dir if not exist.
mkdir -p $ZUNKAPATH/xml

# Last downloaded XML file.
F_LAST=$ZUNKAPATH/xml/aldo/aldo-products.xml
# XML file backup.
F_BACKUP=$ZUNKAPATH/xml/aldo/aldo-products-$(date +%Y-%h-%d-%H%M%S).xml

# Download xml file.
curl "http://webservice.aldo.com.br/asp.net/ferramentas/integracao.ashx?u=146612&p=zunk4c" > $F_BACKUP

# Copy as last xml file.
cp $F_BACKUP $F_LAST

# curl "http://webservice.aldo.com.br/asp.net/ferramentas/saldoproduto.ashx?u=146612&p=zunk4c&codigo=20764-8&qtde=1&emp_filial=1"