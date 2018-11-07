# build stage
FROM golang:1.10.2 AS build-env
WORKDIR /go/src/github.com/turbosonic/event-hub/src
ADD src .
ADD vendor ./vendor
ADD Gopkg* ./
RUN go test -v ./...	
RUN CGO_ENABLED=0 GOOS=linux go build -o eventhub.exe main.go

# final stage
FROM scratch
COPY --from=build-env /go/src/github.com/turbosonic/event-hub/src/eventhub.exe .
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./eventhub.exe"]