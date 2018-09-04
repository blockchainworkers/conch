FROM golang:1.10-alpine3.7
LABEL author="wupenxin  <wupeaking@gamil.com>"


ENV HOME /opt/conch


# 更新软件包
RUN mkdir -p /opt/conch/bin && apk update && apk add curl && apk add git && apk add gcc && \
apk add libc-dev

#    apt-get install -y git && mkdir -p /bitcoin

# 安装守护进程
RUN curl -L https://github.com/Yelp/dumb-init/releases/download/v1.2.1/dumb-init_1.2.1_amd64 -o /usr/local/bin/dumb-init
RUN chmod +x /usr/local/bin/dumb-init

# clone git and build 
RUN mkdir -p /go/src/github.com/blockchainworkers && cd /go/src/github.com/blockchainworkers && git clone https://github.com/blockchainworkers/conch.git conch && \
cd conch && go build -ldflags "-X github.com/blockchainworkers/conch/version.GitCommit=`git rev-parse --short=8 HEAD`"  -o conchd.x ./cmd/conch/main.go && \
mv conchd.x /opt/conch/bin && chmod +x /opt/conch/bin/conchd.x

RUN rm -rf /go/src/github.com/blockchainworkers

VOLUME ["/opt/conch/data"]

WORKDIR /opt/conch/bin

COPY start.sh /opt/conch/bin
RUN chmod +x start.sh

ENTRYPOINT ["dumb-init", "--"]
CMD ["sh", "-x", "/opt/conch/bin/start.sh"]