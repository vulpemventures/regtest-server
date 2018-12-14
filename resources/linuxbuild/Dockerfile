FROM golang:latest
WORKDIR /
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/vulpemventures
WORKDIR /go/src/github.com/vulpemventures
RUN git clone https://github.com/vulpemventures/regtest-server
WORKDIR /go/src/github.com/vulpemventures/regtest-server
RUN chmod +x scripts/install
RUN ./scripts/install
RUN chmod +x ./scripts/build
RUN scripts/build linux amd64
