FROM golang:latest
WORKDIR /
RUN mkdir -p /go/src/github.com/vulpemventures
WORKDIR /go/src/github.com/vulpemventures
RUN git clone https://github.com/vulpemventures/regtest-server
WORKDIR /go/src/github.com/vulpemventures/regtest-server
RUN chmod +x scripts/install
RUN ./scripts/install
RUN chmod +x ./scripts/build
RUN scripts/build linux amd64
