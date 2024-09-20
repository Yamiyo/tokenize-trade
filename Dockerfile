FROM golang:1.22 AS go-builder
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .
RUN go build -o tokenize-trade /app/cmd/

FROM alpine:3.17

RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add git
WORKDIR /tokenize-trade
COPY --from=go-builder /app/tokenize-trade /tokenize-trade/tokenize-trade
COPY --from=go-builder /app/conf.d /tokenize-trade/conf.d
COPY --from=go-builder /app/web/dist /tokenize-trade/web/dist
ENTRYPOINT ["/tokenize-trade/tokenize-trade"]
