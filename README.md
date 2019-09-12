# Pangolin

Pangolin is a Go implenmentation of VPN. 
* Support TCP/UDP/[PTCP](https://github.com/xitongsys/ptcp) connection. (I suggest PTCP, which has the same performance with UDP and avoid some UDP issues)
* For PTCP/TCP, it supports multi-user authentication, encryption transmission. For UDP, no authentication.
* Using Docker, it supports Linux/Windows/Mac.
* For client, it supports Linux/Windows/Android now.

## Server 
Pangolin server can only run natively on Linux. But you can use docker to run it on Windows and Mac.

* Download the latest release package.

* ```cd pangolin_linux/``` 

* Change the environment variables in start.sh to your own.

```bash
SERVERIP=0.0.0.0
SERVERPORT=12345
TOKENS='["token01", "token02"]'
ROLE=SERVER
```

* Start the pangolin server.
```bash
bash start.sh
```

## Client

### Linux
Same steps with the server. But change the ```ROLE``` to ```CLIENT```.
```bash
SERVERIP=your.server.ip.address
SERVERPORT=12345
TOKENS='["token01", "token02"]'
ROLE=CLIENT
```

### Windows
* [Windows Client](https://github.com/xitongsys/pangolin-win)

### Android
* [Android Client](https://github.com/xitongsys/pangolin-android)

## Docker

### Build 
* ```cd pangolin/docker```
* Change the variables in ```pangolin_docker.sh```

```bash
SOURCE_DIR="../"
SERVERIP="0.0.0.0"
SERVERPORT="12345"
TOKENS='["token01", "token02"]'
ROLE="SERVER"
```

* Build ```bash pangolin_docker.sh build```


### Run/Stop
```bash
bash pangolin_docker.sh start
bash pangolin_docker.sh stop
```

## Status
This project is still in progress and you can contribute to it. Anything is welcome !
* Mac/iOS client
* Improve Android client
* Add UT/Doc


