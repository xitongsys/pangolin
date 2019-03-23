#Pangolin

Pangolin is a pure-go implenmentation of TUN VPN.

## Deploy
### Server
1. Run script/server.sh on your VPN server
2. go run main.go -role server -server 0.0.0.0:12345

### Cient
* Linux host
1. Run script/client.sh on your Linux host
2. go run main.go -role client -server \$YOUR_SERVER_ADD

* Android
1. Install apk/pangolin.apk
2. Run it




