FROM ubuntu
ENTRYPOINT bash /pangolin/start.sh
ENV SERVERIP "47.240.40.78"
ENV SERVERPORT "12345"
ENV ROLE "SERVER"

RUN apt update
RUN apt install -y iproute2 iptables net-tools
COPY pangolin /pangolin
RUN chmod 777 /pangolin/main
