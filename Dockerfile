FROM golang:alpine AS builder

RUN apk add --update --no-cache git openssh
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/app

COPY . .

RUN dep ensure && go build -o docker_stats_exporter ./src/main.go



FROM alpine

COPY --from=builder /go/src/app/docker_stats_exporter /bin/docker_stats_exporter

EXPOSE 9100
CMD /bin/docker_stats_exporter