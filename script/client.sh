SERVER=""
GW=""

ip tuntap add dev tun0 mod tun
ip addr add 10.0.0.2/24 dev tun0
ip link set tun0 up

route add $SERVER gw $GW
route add default gw 10.0.0.1

