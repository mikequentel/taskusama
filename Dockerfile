# syntax=docker/dockerfile:1

FROM golang:1.26.0-alpine3.23 AS build

WORKDIR /src

# sets reproducible build flags
ENV GOFLAGS="-trimpath -buildvcs=false"
ENV LDFLAGS="-s -w"

# caches deps
COPY go.mod go.sum ./
RUN go mod download

# copies source
COPY . ./

# builds a static-ish Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="${LDFLAGS}" -o /out/taskusama ./cmd/server


FROM alpine:3.23

RUN apk --no-cache add ca-certificates \
    && addgroup -g 1000 -S taskusama \
    && adduser  -u 1000 -S -G taskusama taskusama

# puts runtime files under /app so relative paths like web/templates/* work
WORKDIR /app

COPY --from=build /out/taskusama /usr/local/bin/taskusama
COPY --from=build /src/web ./web

USER taskusama

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/taskusama"]

