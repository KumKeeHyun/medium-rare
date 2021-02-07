FROM golang:1.13-alpine as builder

WORKDIR /go/src/github.com/KumKeeHyun/medium-rare/trend-service
COPY . .

RUN export GO111MODULE=on
RUN go mod download

RUN go build -o main .

FROM alpine:latest
WORKDIR /bin/
COPY --from=builder /go/src/github.com/KumKeeHyun/medium-rare/trend-service .

ENTRYPOINT [ "./main" ]