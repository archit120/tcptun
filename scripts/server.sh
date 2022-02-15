#!/usr/bin/env bash

# ${gateway IP} ${device name}
ip addr add $1 dev $2
ip link set dev $2 mtu 1400
ip link set dev $2 up
sysctl -w net.ipv4.ip_forward=1
update-alternatives --set iptables /usr/sbin/iptables-legacy
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
