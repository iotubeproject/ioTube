FROM golang:1.17-alpine as build

WORKDIR apps/iotube-witness

RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir -p $GOPATH/pkg/linux_amd64/github.com/iotubeproject/ && \
    make clean build-witness build-explorer build-addtokenpair

FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN mkdir -p /etc/iotube/
COPY --from=build /go/apps/iotube-witness/bin/witness /usr/local/bin/witness
# COPY --from=build /go/apps/iotube-witness/bin/relayer /usr/local/bin/relayer
COPY --from=build /go/apps/iotube-witness/bin/explorer /usr/local/bin/explorer
COPY --from=build /go/apps/iotube-witness/bin/addtokenpair /usr/local/bin/addtokenpair

CMD [ "witness" ]
