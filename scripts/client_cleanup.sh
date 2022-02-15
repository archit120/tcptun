#!/usr/bin/env bash

GW=$1
if [[ $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    :
else
    GW=$(getent hosts $1 | awk 'NR==1{ print $1 }')
fi
ip route del $GW