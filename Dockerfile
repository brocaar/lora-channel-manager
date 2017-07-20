FROM golang:1.8

ENV PROJECT_PATH=/go/src/github.com/brocaar/lora-gateway-config
ENV PATH=$PATH:$PROJECT_PATH/build

# install tools
RUN go get github.com/golang/lint/golint
RUN go get github.com/kisielk/errcheck
RUN go get github.com/smartystreets/goconvey

# setup work directory
RUN mkdir -p $PROJECT_PATH
WORKDIR $PROJECT_PATH
