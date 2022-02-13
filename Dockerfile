FROM golang:1.13

RUN mkdir -p /go/src/bitbucket.org/cloud-platform/vnpt-sso-usermgnt
WORKDIR /go/src/bitbucket.org/cloud-platform/vnpt-sso-usermgnt
ADD . /go/src/bitbucket.org/cloud-platform/vnpt-sso-usermgnt

RUN mkdir -p /go/src/go.opentelemetry.io/otel
RUN cp -r ./build/otel@v0.7.0/. /go/src/go.opentelemetry.io/otel/
RUN mkdir -p /go/src/github.com/go-redis
RUN cp -r ./build/redis /go/src/github.com/go-redis/

#RUN cp -r ./build/confluent-kafka-go-1.4.2 /go/src/github.com/confluent-kafka-go-1.4.2/
#RUN cp -r ./build/confluentinc /go/src/github.com/confluentinc/
RUN go get gopkg.in/confluentinc/confluent-kafka-go.v1/kafka
#Build librdkafka from source 
RUN apt update -y
RUN apt install -y librdkafka-dev
#
RUN go get .
RUN mkdir /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server .
RUN cp -rp ./config.json /app/
RUN rm -r /go/src/bitbucket.org/cloud-platform

EXPOSE 1323
WORKDIR /app
ENTRYPOINT ["/app/server"]
