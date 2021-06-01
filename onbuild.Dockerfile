FROM registry.cn-beijing.aliyuncs.com/g7/env-golang:golang

ENV CGO_ENABLED 0

RUN sed -i "s|http://dl-cdn.alpinelinux.org|http://mirrors.aliyun.com|g" /etc/apk/repositories

RUN echo "2018-06"

RUN apk add --no-cache curl git openssh wget unzip

RUN git config --global url."https://gitlab-ci-token:AMUrtHnxw4RQ8wJ6nzHD@git.chinawayltd.com/".insteadOf "https://git.chinawayltd.com/"
RUN git config --global url."https://morlay:d2cd29e12d833ce737fd3a49590c8872b68dee62@github.com/".insteadOf "https://github.com/"
RUN git config --global url."https://morlay:d2cd29e12d833ce737fd3a49590c8872b68dee62@github.com/golang/".insteadOf "https://go.googlesource.com/"

RUN go get -u golang.org/x/vgo \
    && go get -u github.com/kardianos/govendor \
    && go get -u github.com/morlay/gin-swagger

RUN ln -s /tmp/gomodules /go/src/mod

COPY . /go/src/git.chinawayltd.com/golib/tools
RUN cd /go/src/git.chinawayltd.com/golib/tools && go install
