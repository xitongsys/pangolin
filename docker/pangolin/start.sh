function start_server
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.2/8 dev tun0
	ip link set tun0 up
	ip=`ip addr show dev "eth0" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to-source $ip
	
	sed -i "s/SERVERPORT/$SERVERPORT/g" /pangolin/cfg.json
	/pangolin/main -c /pangolin/cfg.json 
}

function start_client
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.2/8 dev tun0
	ip link set tun0 up
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o tun0 -j SNAT --to-source 10.0.0.2
	route add default gw 10.0.0.1
	
	gw=`route -n | awk '$1 == "0.0.0.0" {print $2}'`
	route add $SERVERIP gw $gw
	sed -i "s/SERVERIP/$SERVERIP/g" /pangolin/cfg_client.json
	sed -i "s/SERVERPORT/$SERVERPORT/g" /pangolin/cfg_client.json
	/pangolin/main -c /pangolin/cfg_client.json 
}


[[ "$ROLE" -eq "SERVER" ]] && start_server
[[ "$ROLE" -eq "CLIENT" ]] && start_client
