SERVER=""
DEVICE=""

ip tuntap add dev tun0 mod tun
ip addr add 10.0.0.2/24 dev tun0
ip link set tun0 up


cat 1 >> /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -o $DEVICE -j SNAT --to-source $SERVER
