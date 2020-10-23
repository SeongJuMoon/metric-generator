from golang:1.14 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN GOOS=linux GO111MOUDLE=on CGO_ENABLED=0 go build -a -o main . 

from alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=builder /workspace/main .
ENTRYPOINT ["/main"] 