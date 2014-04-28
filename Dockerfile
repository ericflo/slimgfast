FROM ubuntu:14.04

MAINTAINER Eric Florenzano <floguy@gmail.com>

RUN echo "deb http://archive.ubuntu.com/ubuntu trusty main universe" > /etc/apt/sources.list
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y curl git bzr mercurial

WORKDIR /tmp

RUN curl -O https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz
RUN tar xvfz godeb-amd64.tar.gz
RUN ./godeb install

RUN adduser --home /home/slimgfast slimgfast

WORKDIR /home/slimgfast

RUN mkdir -p go/src/github.com/ericflo
RUN mkdir go/bin
RUN mkdir go/pkg
RUN chown slimgfast:slimgfast -R .

ENV GOPATH /home/slimgfast/go

RUN go get github.com/ericflo/slimgfast/slimgfastd

USER slimgfast

ENTRYPOINT ["/home/slimgfast/go/bin/slimgfastd"]

EXPOSE 4400