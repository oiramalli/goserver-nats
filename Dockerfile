FROM golang
RUN go get github.com/nats-io/nats.go/
RUN useradd -ms /bin/bash goserveruser
RUN mkdir /goserver
WORKDIR /goserver
ADD /main/ /goserver/
RUN go build -o main .
USER goserveruser
EXPOSE 8080
CMD ["./main"]
