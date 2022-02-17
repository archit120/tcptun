#!/usr/bin/env bash

# ${client IP} ${device name} ${server ip} ${server gateway}
sleep 1

dhclient $1 -v
L1=$(ip r | grep $1 | awk '{print $1; exit}')
IFS='/' read -r -a array <<< "${L1}"
IPA=${array[0]}
IFS='.' read -r -a array <<< "${IPA}"
IPA=${array[0]}.${array[1]}.${array[2]}.1
ip route add 0.0.0.0/1 via ${IPA} dev $1 onlink
ip route add 128.0.0.0/1 via ${IPA} dev $1 onlink
