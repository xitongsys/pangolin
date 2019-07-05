FROM ubuntu
USER root
ENTRYPOINT bash /pangolin/start.sh
ENV SERVERIP "{SERVERIP}"
ENV ROLE "SERVER"

RUN apt update
RUN apt install -y iproute2 iptables net-tools dos2unix
COPY pangolin /pangolin
RUN chmod 777 /pangolin/main
RUN dos2unix /pangolin/start.sh
