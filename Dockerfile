FROM starudream/golang AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on go build -o clash-speedtest .

RUN apk add --no-cache upx && upx clash-speedtest

FROM starudream/alpine-glibc:latest

COPY --from=builder /build/clash-speedtest /clash-speedtest

WORKDIR /

CMD /clash-speedtest
