ip tuntap add dev tun0 mod tun
ip addr add 10.0.0.2/24 dev tun0
ip link set tun0 up


cat 1 >> /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -o ens3 -j SNAT --to-source 139.180.132.42
