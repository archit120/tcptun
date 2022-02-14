# tcptun

The goal of this project is to create a simple internet tunnel to bypass firewalls that only allow traffic on TCP port 443. 


Commands to run on client are - 
```
ip addr add 10.1.0.2/24 dev tun0
ip link set dev tun0 mtu 1400
ip link set dev tun0 up
ip route add 0.0.0.0/1 dev tun0
ip route add 128.0.0.0/1 dev tun0
IP=$(getent hosts server | awk '{ print $1 }')
ip route add ${IP}/32 dev eth0
```

Commands to run on server are -

```
ip addr add 10.1.0.1/24 dev tun0
ip link set dev tun0 mtu 1400
ip link set dev tun0 up
sysctl -w net.ipv4.ip_forward=1
update-alternatives --set iptables /usr/sbin/iptables-legacy
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

```