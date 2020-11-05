FROM golang:1.15-alpine as builder

RUN apk update && apk upgrade && \
apk add --no-cache bash git openssh
RUN apk add --no-cache autoconf automake libtool gettext gettext-dev make g++ texinfo curl ffmpeg


# Add Maintainer Info
LABEL maintainer="<PaulBarrie <barrie_p@etna-alternance.net>"
ENV GO111MODULE=on

COPY main.go /go/src/api/
WORKDIR /go/src/api
COPY . /usr/local/go/src

COPY go.mod go.sum ./ 

# Vendoring dependencies
RUN go mod tidy
RUN go mod vendor
COPY . .
# RUN mv config controllers docs models router vendor
# Build the Go api


# WORKDIR /root
# RUN wget https://github.com/emcrisostomo/fswatch/releases/download/1.14.0/fswatch-1.14.0.tar.gz
# RUN tar -xvzf fswatch-1.14.0.tar.gz
# WORKDIR /root/fswatch-1.14.0
# RUN ./configure
# RUN make
# RUN make install

ADD  . /go/src/api
#RUN mv config controllers docs models router vendor

EXPOSE 3001
