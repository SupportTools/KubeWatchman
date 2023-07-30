FROM golang:1.20 AS builder

WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go test -v ./...
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/KubeWatchman .

FROM scratch
COPY --from=builder /go/bin/KubeWatchman /KubeWatchman
CMD ["/KubeWatchman"]