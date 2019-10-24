#!/usr/bin/env bash

# $1 - table
# $2 - key
# $3 - file

if [[ -z $3 ]]; then
    echo Usage: $0 table key file
    exit 1
fi

# count=1
while read line; do
    # Find table.
    if [[ $line =~ \[$1\] ]]; then
        # echo $line
        # Find key.
        while read line; do
            if [[ $line =~ $2 ]]; then
                # echo $line
                IFS='=' read -ra AR <<< $line
                val=${AR[1]}
                # Replace first ${val/substring/replacement}
                # Replace all ${val//substring/replacement}
                echo ${val//\"}
                exit 0
            fi
        done
    fi
    # echo $count
    # let count=count+1
done < $3

# for line in $(cat $1); do
    # if [[ $line =~ \[zunkasrv\] ]]; then
        # echo $line
    # fi
# done
