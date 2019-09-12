SOURCE_DIR="../"
SERVERIP="0.0.0.0"
SERVERPORT="12345"
TOKENS='["token01", "token02"]'
ROLE="SERVER"

function build () {
	go build -o pangolin/main $SOURCE_DIR
	docker build -t pangolin .
}

function start () {
    docker run --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=$ROLE --env SERVERIP=$SERVERIP --env SERVERPORT=$SERVERPORT --env TOKENS="$TOKENS" pangolin
}

function stop() {
	docker ps | grep pangolin | awk '{print $1}' | xargs -I {} docker kill {} 
	[[ "$ROLE" == "CLIENT" ]] && route del default
}


cmd=$1
[[ "$cmd" == "build" ]] && build
[[ "$cmd" == "start" ]] && start
[[ "$cmd" == "stop" ]] && stop
