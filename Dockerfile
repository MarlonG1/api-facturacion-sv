FROM alpine:latest
WORKDIR /app
COPY dte-service .
COPY config ./config
EXPOSE 7319
CMD ["./dte-service"]
