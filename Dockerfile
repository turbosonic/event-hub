# build stage
FROM golang:1.10.2 AS build-env
RUN go get github.com/eclipse/paho.mqtt.golang
RUN go get pack.ag/amqp
RUN go get github.com/joho/godotenv
RUN go get gopkg.in/yaml.v2
RUN go get gopkg.in/olivere/elastic.v5
RUN go get github.com/satori/go.uuid
RUN go get -u github.com/Azure/azure-service-bus-go/...
RUN go get github.com/Microsoft/ApplicationInsights-Go/appinsights
WORKDIR /go/src/contino.visualstudio.com/event-hub/src
ADD src .
RUN go test -v ./...	
RUN CGO_ENABLED=0 GOOS=linux go build -o eventhub.exe main.go

# final stage
FROM scratch
COPY --from=build-env /go/src/contino.visualstudio.com/event-hub/src/eventhub.exe .
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./eventhub.exe"]