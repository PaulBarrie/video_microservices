FROM alpine:latest

USER root

RUN addgroup -S docker &&\
    adduser -S docker $USER

ENV GOROOT /usr/lib/go 
ENV GOPATH /go 
ENV PATH "/opt/gtk/bin:${PATH}"
ENV GO111MODULE on
WORKDIR /root
RUN wget https://github.com/emcrisostomo/fswatch/releases/download/1.14.0/fswatch-1.14.0.tar.gz &&\
    tar -xvzf fswatch-1.14.0.tar.gz

WORKDIR /root/fswatch-1.14.0
RUN ./configure && make && make install

# Configure Go

RUN mkdir -p ${GOPATH}/src/cve-report ${GOPATH}/bin/cv-report

WORKDIR ${GOPATH}/src

COPY . .
# RUN go install cve-report
RUN go mod tidy &&\
    go mod vendor

ENTRYPOINT ["make"]
CMD [ "serve" ]