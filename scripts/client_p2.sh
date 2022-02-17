#!/usr/bin/env bash

# ${client IP} ${device name} ${server ip} ${server gateway}
sleep 1
ip link set dev $1 mtu 1400
ip link set dev $1 up
dhclient $1 -v
L1=$(ip r | grep $1 | grep -v default | awk '{print $1; exit}')
IFS='/' read -r -a array <<< "${L1}"
IPA=${array[0]}
IFS='.' read -r -a array <<< "${IPA}"
IPA=${array[0]}.${array[1]}.${array[2]}.1
ip route add 0.0.0.0/1 via ${IPA} dev $1
ip route add 128.0.0.0/1 via ${IPA} dev $1 
