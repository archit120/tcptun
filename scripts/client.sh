#!/usr/bin/env bash

# ${client IP} ${device name} ${server ip}
ip addr add $1 dev $2
ip link set dev $2 mtu 1400
ip link set dev $2 up
ip route add 0.0.0.0/1 dev $2
ip route add 128.0.0.0/1 dev $2
IP=$(ip route | awk '/default/ {print $3; exit}')
ip route add $3/32 via ${IP} dev eth0
