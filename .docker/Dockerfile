FROM golang:1.23-alpine AS builder

WORKDIR /build
RUN apk update --no-cache && apk add --no-cache tzdata
COPY ../../ .
RUN go build -o api cmd/server/main.go
RUN go build -o seeds cmd/seeds/main.go

FROM alpine:3.14 AS production
RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ Europe/Moscow
WORKDIR /build
COPY --from=builder build/api /build/api
COPY --from=builder build/seeds /build/seeds
CMD ["./api"]

