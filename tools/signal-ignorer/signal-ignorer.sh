#!/usr/bin/env bash

function signal_handler()
{
    echo "singal ignored"
}

if [ $# -gt 0 ];
    then trap signal_handler "$1"
fi

echo $$
while true
do
    sleep 0.1
done
