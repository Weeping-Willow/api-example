FROM golang:1.17 AS builder

RUN apt-get update \ 
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends make

WORKDIR /go/src/api-example
COPY go.* ./

RUN go env -w CGO_ENABLED=0 && \
    go mod download

COPY ./ ./
RUN make build

FROM alpine:latest  

WORKDIR /root/
COPY --from=builder /go/src/api-example/bin ./
RUN chmod +x ./api-example

ENV PORT=80
ENV DATABASE_URL=mongodb://admin:ASjkjasd13123@localhost:27017/?authSource=example

CMD ["./api-example"]