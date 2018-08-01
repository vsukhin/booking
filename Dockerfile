FROM golang:alpine

ADD . /go/src/github.com/vsukhin/booking

WORKDIR /go/src/github.com/vsukhin/booking

RUN go install .

EXPOSE 3000

RUN addgroup booking && adduser -S -G booking booking

USER booking

ENTRYPOINT ["/go/bin/booking"]
