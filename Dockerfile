FROM golang:latest

LABEL maintainer="Lucas Golino <contato@golino.space>"

RUN mkdir /opt/cloudflare-zone-manager
ADD . /opt/cloudflare-zone-manager
WORKDIR /opt/cloudflare-zone-manager

ENV GO111MODULE=on

RUN go mod download

RUN cd /opt/cloudflare-zone-manager/cmd && go build -buildmode=pie -o czm
RUN cd /opt/cloudflare-zone-manager/modules/ && ./build-modules.sh

CMD ["./cmd/czm"]