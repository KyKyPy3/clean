# Builder
FROM golang:1.21.3-alpine3.18 as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app && \
    mkdir /app/config

WORKDIR /app

COPY --from=builder /app/clean /app/
COPY --from=builder /app/config/config-docker.yml /app/config

ENV config=config/config-docker.yml

EXPOSE 8080

CMD /app/clean