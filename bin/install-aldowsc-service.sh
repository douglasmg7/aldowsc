#! /usr/bin/env bash

# # Should run as a root.
# if [ "$EUID" -ne 0 ]; then 
  # echo "Please run as root"
  # exit
# fi

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 

# GOPATH must be defined.
[[ -z "$GOPATH" ]] && printf "error: GOPATH enviorment not defined.\n" >&2 && exit 1 

# GS must be defined.
[[ -z "$GS" ]] && printf "error: GS enviorment not defined.\n" >&2 && exit 1 

# Script not exist.
[[ ! -f $GS/aldowsc/bin/fetch-xml-products-and-process.sh ]] && printf "error: script $GS/aldowsc/bin/fetch-xml-products-and-process.sh not exist.\n" >&2 && exit 1 

# Uninstall script not exist.
[[ ! -f $GS/aldowsc/bin/uninstall-aldowsc-service.sh ]] && printf "error: script $GS/aldowsc/bin/uninstall-aldowsc-service.sh not exist.\n" >&2 && exit 1

# Remove aldowsc timer and aldo service.
./uninstall-aldowsc-service.sh

# Make aldowsc script wide system accessible.
echo Creating symobolic link for aldowsc script...
sudo ln -s $GOPATH/bin/aldowsc /usr/local/bin/aldowsc

# Create aldo timer.
echo "creating '/lib/systemd/system/aldowsc.timer'..."
sudo bash -c 'cat << EOF > /lib/systemd/system/aldowsc.timer
[Unit]
Description=aldowsc timer

[Timer]
OnCalendar=*-*-* 00:00:01
# OnCalendar=*-*-* 01:00:00
# OnCalendar=*-*-* 02:00:00
# OnCalendar=*-*-* 03:00:00
# OnCalendar=*-*-* 04:00:00
# OnCalendar=*-*-* 05:00:00
# OnCalendar=*-*-* 06:00:00
# OnCalendar=*-*-* 07:00:00
# OnCalendar=*-*-* 08:00:00
# OnCalendar=*-*-* 09:00:00
# OnCalendar=*-*-* 10:00:00
# OnCalendar=*-*-* 11:00:00
# OnCalendar=*-*-* 12:00:00
# OnCalendar=*-*-* 13:00:00
# OnCalendar=*-*-* 14:00:00
# OnCalendar=*-*-* 15:00:00
# OnCalendar=*-*-* 16:00:00
# OnCalendar=*-*-* 17:00:00
# OnCalendar=*-*-* 18:00:00
# OnCalendar=*-*-* 19:00:00
# OnCalendar=*-*-* 20:00:00
# OnCalendar=*-*-* 21:00:00
# OnCalendar=*-*-* 22:00:00
# OnCalendar=*-*-* 23:00:00

Persistent=true

[Install]
WantedBy=timers.target
EOF'

# Create aldo service.
echo "creating '/lib/systemd/system/aldowsc.service'..."
sudo GS=$GS ZUNKAPATH=$ZUNKAPATH ZUNKA_ALDOWSC_DB=$ZUNKA_ALDOWSC_DB bash -c 'cat << EOF > /lib/systemd/system/aldowsc.service
[Unit]
Description=aldowsc

[Service]
Type=oneshot
User=douglasmg7
Environment="ZUNKAPATH=$ZUNKAPATH"
Environment="ZUNKA_ALDOWSC_DB=$ZUNKA_ALDOWSC_DB"
ExecStart=$GS/aldowsc/bin/fetch-xml-products-and-process.sh
EOF'

sudo systemctl start aldowsc.timer
sudo systemctl enable aldowsc.timer
