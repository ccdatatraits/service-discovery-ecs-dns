FROM golang:1.6.3

RUN go get go get github.com/gorilla/mux

COPY src /go/src/github.com/mediamath/techsol/log-streamer

RUN go install github.com/mediamath/techsol/log-streamer

EXPOSE 8081

CMD ["/go/bin/log-streamer"]
