#!/usr/bin/env bash

# ${client IP} ${device name} ${server ip} ${server gateway}
ip addr add $1 dev $2
ip link set dev $2 mtu 1400
ip link set dev $2 up
ip route add 0.0.0.0/1 via $4 dev $2
ip route add 128.0.0.0/1 via $4 dev $2
IP=$(ip route | awk '/default/ {print $3; exit}')
DEV=$(ip route | awk '/default/ {print $5; exit}')
GW=$3
if [[ $3 =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    :
else
    GW=$(getent hosts $3 | awk 'NR==1{ print $1 }')
fi
echo "ip route add $GW/32 via ${IP} dev ${DEV}"
