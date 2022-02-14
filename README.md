# tcptun

The goal of this project is to create a simple internet tunnel to bypass firewalls that only allow traffic on TCP port 443. 


Commands to run on client are - 
```
ip addr add 10.1.0.2/24 dev tun0
ip route add 0.0.0.0/1 dev tun0
ip route add 128.0.0.0/1 dev tun0
ip link set dev tun0 up
ping 10.1.0.254
```

Commands to run on server are -

```
sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -o eth1 -j MASQUERADE
```