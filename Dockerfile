FROM golang:1-alpine AS builder

# default timezon
# override it with `--build-arg TIMEZONE=xxxx`
ARG TIMEZONE=Asia/Shanghai

ENV TZ ${TIMEZONE}

COPY . /go/src/github.com/mritd/notibot

WORKDIR /go/src/github.com/mritd/notibot

RUN set -ex \
    && apk add gcc musl-dev git tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && export version=$(git describe --tags --always) \
    && export build_date=$(date '+%F %T') \
    && export commit_hash=$(git rev-parse HEAD) \
    && go install -trimpath -ldflags "-w -s \
        -X \"main.version=${version}\" \
        -X \"main.build=${build_date}\" \
        -X \"main.commit=${commit_hash}\""

FROM alpine

LABEL maintainer="mritd <mritd@linux.com>"
LABEL org.opencontainers.image.source="https://github.com/mritd/notibot"
LABEL org.opencontainers.image.description="Telegram Notification Bot"
LABEL org.opencontainers.image.licenses="Apache-2.0"

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

# default timezon
# override it with `--build-arg TIMEZONE=xxxx`
ARG TIMEZONE=Asia/Shanghai

ENV TZ ${TIMEZONE}

RUN set -ex \
    && apk upgrade \
    && apk add --no-cache bash ca-certificates tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

COPY --from=builder /go/bin/notibot /usr/local/bin/notibot

EXPOSE 8080

CMD ["notibot"]