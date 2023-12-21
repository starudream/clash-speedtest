FROM starudream/alpine

WORKDIR /

COPY clash-speedtest /clash-speedtest

CMD /clash-speedtest
