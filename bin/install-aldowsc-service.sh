#! /usr/bin/env bash

# Should run as a root.
# if [ "$EUID" -ne 0 ]; then 
  # echo "Please run as root"
  # exit
# fi

# Script not exist.
[[ ! -f $GS/aldowsc/bin/fetch-xml-products-and-process.sh ]] && printf "error: script $GS/aldowsc/bin/fetch-xml-products-and-process.sh not exist.\n" >&2 && exit 1 

# Uninstall script not exist.
[[ ! -f $GS/aldowsc/bin/uninstall-aldowsc-service.sh ]] && printf "error: script $GS/aldowsc/bin/uninstall-aldowsc-service.sh not exist.\n" >&2 && exit 1


# Remove aldowsc timer and aldo service.
./uninstall-aldowsc-service.sh

# Create aldo timer.
echo "creating '/lib/systemd/system/aldowsc.timer'"
sudo bash -c 'cat << EOF > /lib/systemd/system/aldowsc.timer
[Unit]
Description=aldowsc timer

[Timer]
OnCalendar=*-*-* 00:00:00
OnCalendar=*-*-* 01:00:00
OnCalendar=*-*-* 02:00:00
OnCalendar=*-*-* 03:00:00
OnCalendar=*-*-* 04:00:00
OnCalendar=*-*-* 05:00:00
OnCalendar=*-*-* 06:00:00
OnCalendar=*-*-* 07:00:00
OnCalendar=*-*-* 08:00:00
OnCalendar=*-*-* 09:00:00
OnCalendar=*-*-* 10:00:00
OnCalendar=*-*-* 11:00:00
OnCalendar=*-*-* 12:00:00
OnCalendar=*-*-* 13:00:00
OnCalendar=*-*-* 14:00:00
OnCalendar=*-*-* 15:00:00
OnCalendar=*-*-* 16:00:00
OnCalendar=*-*-* 17:00:00
OnCalendar=*-*-* 18:00:00
OnCalendar=*-*-* 19:00:00
OnCalendar=*-*-* 20:00:00
OnCalendar=*-*-* 21:00:00
OnCalendar=*-*-* 22:00:00
OnCalendar=*-*-* 23:00:00

Persistent=true

[Install]
WantedBy=timers.target
EOF'

# Create aldo service.
echo "creating '/lib/systemd/system/aldowsc.service'"
sudo GS=$GS bash -c 'cat << EOF > /lib/systemd/system/aldowsc.service
[Unit]
Description=aldowsc

[Service]
Type=oneshot
User=douglasmg7
ExecStart=$GS/aldowsc/bin/fetch-xml-products-and-process.sh
EOF'

systemctl start aldowsc.timer
systemctl enable aldowsc.timer
