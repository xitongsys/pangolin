#SERVERIP/SERVERPORT/TOKENS
function start_server ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.2/8 dev tun0
	ip link set tun0 up
	ip link set dev tun0 mtu 1200
	ip=`ip addr show dev "eth0" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to-source $ip
	iptables -P FORWARD ACCEPT
	/pangolin/main -c /pangolin/configs/cfg_server.json -l debug
}

function start_client ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.22/8 dev tun0
	ip link set tun0 up
	ip link set dev tun0 mtu 1200
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o tun0 -j SNAT --to-source 10.0.0.22
	iptables -P FORWARD ACCEPT
	
	gw=`route -n | awk '$1 == "0.0.0.0" {print $2}'`
	route add $SERVERIP gw $gw
	route add default gw 10.0.0.1
	echo "nameserver 8.8.8.8" > /etc/resolv.conf
	/pangolin/main -c /pangolin/configs/cfg_client.json  -l debug
}

function replace ()
{	
	sed -i "s/{SERVERIP}/$SERVERIP/g" $1
	sed -i "s/{SERVERPORT}/$SERVERPORT/g" $1
	sed -i "s/{TOKENS}/$TOKENS/g" $1
}

replace /pangolin/configs/cfg_client.json
replace /pangolin/configs/cfg_server.json


[[ "$ROLE" == "SERVER" ]] && start_server
[[ "$ROLE" == "CLIENT" ]] && start_client

echo "ERROR: pangolin exit"
