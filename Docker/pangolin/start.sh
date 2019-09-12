#SERVERIP=0.0.0.0
#SERVERPORT=12345
#TOKENS='["token01", "token02"]'
#ROLE=CLIENT

function start_server ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.2/8 dev tun0
	ip link set tun0 up
	ip link set dev tun0 mtu 1400
	ifc=`route -n | awk '{if($1=="0.0.0.0") print $8}'`
	INTERFACE="$ifc"
	ip=`ip addr show dev "$ifc" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`
	SERVERIP=$ip
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o "$ifc" -j SNAT --to-source $ip
	iptables -P FORWARD ACCEPT
	iptables -A INPUT -p tcp --destination-port `expr $SERVERPORT + 1` -j DROP

	replace configs/cfg_server.json.template > configs/cfg_server.json
	./main -c configs/cfg_server.json 
}

function start_client ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.22/8 dev tun0
	ip link set tun0 up
	ip link set dev tun0 mtu 1400
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o tun0 -j SNAT --to-source 10.0.0.22
	iptables -P FORWARD ACCEPT
	ifc=`route -n | awk '{if($1=="0.0.0.0") print $8}'`
	INTERFACE="$ifc"
	
	gw=`route -n | awk '$1 == "0.0.0.0" {print $2}'`
	route add $SERVERIP gw $gw
	route add default gw 10.0.0.1
	echo "nameserver 8.8.8.8" > /etc/resolv.conf

	replace /pangolin/configs/cfg_client.json.template > configs/cfg_client.json
	./main -c ./configs/cfg_client.json 
}


function replace ()
{	
	sed -n "s/{SERVERIP}/$SERVERIP/g" $1
	sed -n "s/{SERVERPORT}/$SERVERPORT/g" $1
	sed -n "s/{TOKENS}/$TOKENS/g" $1
	sed -n "s/{INTERFACE}/$INTERFACE/g" $1
}



[[ "$ROLE" == "SERVER" ]] && start_server
[[ "$ROLE" == "CLIENT" ]] && start_client

echo "ERROR: pangolin exit"
