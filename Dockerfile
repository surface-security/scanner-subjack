FROM golang:1.13.1-alpine3.10 as builder

RUN apk --no-cache add git

ENV CGO_ENABLED=0

WORKDIR /go/src/
ADD go.mod go.sum /go/src/
RUN go mod download

ADD main.go /go/src/
RUN go build -o /subjack -ldflags="-s -w" .


FROM alpine:3.10
RUN apk add --no-cache ca-certificates
# in case we want to save about 5MB...
# FROM scratch

COPY --from=builder /subjack /subjack
ADD fingerprints.json /fingerprints.json

ENTRYPOINT ["/subjack"]
