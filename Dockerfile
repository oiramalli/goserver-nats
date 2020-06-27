FROM golang:alpine
RUN go get github.com/nats-io/nats.go/
RUN mkdir /goserver
ADD /main/ /goserver/
WORKDIR /goserver
RUN go build -o main .
RUN adduser -S -D -H -h /goserver goserveruser
USER goserveruser
EXPOSE 8080
CMD ["./main"]