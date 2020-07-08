FROM golang:1.14-buster

WORKDIR /go
RUN mkdir -p /go/src/github.com/jlandowner

WORKDIR /go/src/github.com/jlandowner
RUN git clone https://github.com/jlandowner/psp-util.git

WORKDIR /go/src/github.com/jlandowner/psp-util
RUN make build
CMD [ "./bin/psp-util" ]