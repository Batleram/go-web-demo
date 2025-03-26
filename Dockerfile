FROM golang:1.24-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go tool templ generate
RUN go build -x -o app

FROM debian:bookworm-slim

WORKDIR /srv

COPY --from=builder /build/app .
COPY migrations /srv/migrations/

CMD [ "/srv/app" ]
