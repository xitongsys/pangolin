SERVERIP="0.0.0.0"
SERVERPORT="12345"
TOKENS='\[\"token01\", \"token02\"\]'
ROLE="CLIENT"


function install () {
	docker build -t pangolin .
}

function start () {
    docker run --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=$ROLE --env SERVERIP=$SERVERIP --env SERVERPORT=$SERVERPORT --env TOKENS="$TOKENS" pangolin
}

function stop() {
	docker 
}


start
