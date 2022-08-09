FROM golang:1.18.3-alpine3.16 as build

ENV GO111MODULE=on

WORKDIR /go/src/github.com/laupse/twitter-analytics-exporter
COPY . .

RUN go build -o /go/bin/twitter-analytics-exporter

FROM alpine:3.16

COPY --from=build /go/bin/twitter-analytics-exporter /go/bin/

CMD [ "/go/bin/twitter-analytics-exporter" ]