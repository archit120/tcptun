#!/usr/bin/env bash


IP=$(ip route | awk '/default/ {print $3; exit}')
DEV=$(ip route | awk '/default/ {print $5; exit}')
GW=$1
if [[ $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    :
else
    GW=$(getent hosts $1 | awk 'NR==1{ print $1 }')
fi
ip route add $GW/32 via ${IP} dev ${DEV}
