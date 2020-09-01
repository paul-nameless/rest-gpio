########################################
FROM golang:1.14.7-alpine3.12 AS builder

WORKDIR /go/src/home/
COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/main -mod=vendor

############
FROM scratch
COPY --from=builder /go/bin/main /main
CMD ["/main"]