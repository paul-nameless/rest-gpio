########################################
FROM golang:1.14.7-alpine3.12 AS builder

RUN apk update && apk add upx

WORKDIR /go/src/home/
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/main -mod=vendor
RUN upx /go/bin/main

############
FROM scratch
COPY --from=builder /go/bin/main /main
CMD ["/main"]