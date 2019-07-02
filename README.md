# Pangolin

Pangolin is a Go implenmentation of TUN VPN. 
* Support Tcp/Udp connection.
* For Tcp, it supports multi-user authentication, encryption transmission. For Udp, no authentication.
* Using Docker, it supports Linux/Windows/Mac.
* For client, it supports Linux/Windows/Mac/Android.

## Deploy
The pangolin can only run natively on Linux. But you can use docker to run it on Windows and Mac.

### Docker
#### Linux
* Go to the docker directory and change the config files at ```docker/pangolin/configs```.

For ```cfg_client.json```, you need to change the ```server``` to your own server ip and add your ```tokens```.
```
docker\pangolin\configs> cat .\cfg_client.json  
{                                 
    "role": "client",
    "server": "147.140.40.78:12345",
    "tun": "10.0.0.22/8",
    "tunname": "tun0",
    "dns": "8.8.8.8",
    "mtu": 1500,
    "protocol": "tcp",
    "tokens": ["token01", "token02"]
}
```

For ```cfg_server.json```, your just need add your ```tokens```.
```
docker\pangolin\configs> cat .\cfg_server.json                                   
{
    "role": "server",
    "server": "0.0.0.0:12345",
    "tun": "10.0.0.2/8",
    "tunname": "tun0",
    "dns": "8.8.8.8",
    "mtu": 1500,
    "protocol": "tcp",
    "tokens": ["token01", "token02"]
}
```

* Go to the docker directory and build your own docker image.
```
docker build -t pangolin .
```

* Run your docker image on server and client sparately. 

On server:
```
docker run --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=SERVER pangolin
```

On Client (similar with server, but your should assign another environmental variable ```SERVERIP``` to your server ip):
```
docker run --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=CLIENT --env SERVERIP=137.140.40.78 pangolin
```

#### Windows
* Follow this [link](https://docs.docker.com/machine/drivers/hyper-v/#2-set-up-a-new-external-network-switch-optional) to create an external VMSwitch and a new docker machine.
* Use the new docker machine and follow the steps of Linux.
* Change your windows host default gateway to the docker
```
route delete 0.0.0.0
route add 0.0.0.0 mask 0.0.0.0 192.168.0.13
#192.168.0.13 is your docker public ip.
```

#### Mac
* Not test yet, but I think it's ok. Maybe you can help :)

#### Android Client
* [Android Client](https://github.com/xitongsys/pangolin-android)

#### iOS Client
* Not supported yet, maybe you can help :)

### Native
Pangolin can only run natively on Linux host. Before running pangolin, you need to setup environment. You can find details in the ```start.sh``` in docker directory
```bash
function start_server ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.2/8 dev tun0
	ip link set tun0 up
	ip=`ip addr show dev "eth0" | awk '$1 == "inet" { sub("/.*", "", $2); print $2 }'`
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to-source $ip
	iptables -P FORWARD ACCEPT
	/pangolin/main -c /pangolin/configs/cfg_server.json 
}

function start_client ()
{
	ip tuntap add dev tun0 mod tun
	ip addr add 10.0.0.22/8 dev tun0
	ip link set tun0 up
	iptables -t nat -F
	iptables -t nat -A POSTROUTING -o tun0 -j SNAT --to-source 10.0.0.22
	iptables -P FORWARD ACCEPT
	
	gw=`route -n | awk '$1 == "0.0.0.0" {print $2}'`
	route add $SERVERIP gw $gw
	route add default gw 10.0.0.1
	echo "nameserver 8.8.8.8" > /etc/resolv.conf
	/pangolin/main -c /pangolin/configs/cfg_client.json 
}


[[ "$ROLE" == "SERVER" ]] && start_server
[[ "$ROLE" == "CLIENT" ]] && start_client

echo "pangolin exit"

tail -f
```




