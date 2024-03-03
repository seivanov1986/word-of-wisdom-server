FROM golang:latest

COPY ./ ./src

WORKDIR "/go/src"

RUN go build -o /go/bin/client ./cmd/client/main.go
RUN go build -o /go/bin/server ./cmd/server/main.go

CMD ["/go/bin/server"]
