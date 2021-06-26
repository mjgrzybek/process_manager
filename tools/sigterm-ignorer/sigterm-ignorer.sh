#!/usr/bin/env bash


function sigterm_handler()
{
    echo "sigterm ignored"
}

trap sigterm_handler SIGTERM

echo $$
while true
do
    sleep 0.1
done
