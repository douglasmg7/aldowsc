#! /usr/bin/env bash

# Remove aldowsc timer and aldo service.
if systemctl list-units --full --all | grep -Fq "aldowsc.timer"; then
    sudo systemctl stop aldowsc.timer 
    sudo systemctl disable aldowsc.timer 
    sudo rm -v /lib/systemd/system/aldowsc.timer
    sudo rm -v /lib/systemd/system/aldowsc.service
    sudo systemctl daemon-reload
    sudo systemctl reset-failed
fi
