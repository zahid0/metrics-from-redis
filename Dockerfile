FROM golang
COPY . /go/src/metrics-from-redis
WORKDIR /go/src/metrics-from-redis
RUN go build

FROM gcr.io/distroless/base:latest
COPY --from=0 /go/src/metrics-from-redis/metrics-from-redis /metrics-from-redis
CMD ["/metrics-from-redis"]
