#! /usr/bin/env bash

# Remove link for all user access.
if [[ -f /usr/local/bin/aldowsc ]]; then
    echo Removing link for aldowsc...
    sudo rm /usr/local/bin/aldowsc
fi

# Remove aldowsc timer and aldo service.
if systemctl list-units --full --all | grep -Fq "aldowsc.timer"; then
    echo Removing service...
    sudo systemctl stop aldowsc.timer 
    sudo systemctl disable aldowsc.timer 
    sudo rm -v /lib/systemd/system/aldowsc.timer
    sudo rm -v /lib/systemd/system/aldowsc.service
    sudo systemctl daemon-reload
    sudo systemctl reset-failed
fi
