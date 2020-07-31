FROM golang:1.14 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make all

FROM golang:alpine3.11

WORKDIR /usr/share/zoneinfo
COPY --from=builder /usr/share/zoneinfo .

WORKDIR /app
ENV TZ=UTC

COPY --from=builder /app/cache_service .
COPY --from=builder /app/docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["./docker-entrypoint.sh"]
