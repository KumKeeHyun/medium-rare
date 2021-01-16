FROM golang:1.13-alpine

WORKDIR $GOPATH/src/github.com/KumKeeHyun/medium-rare/user-service
COPY . $GOPATH/src/github.com/KumKeeHyun/medium-rare/user-service

RUN export GO111MODULE=on
RUN go mod download

RUN go build -o main .

ENTRYPOINT [ "./main" ]