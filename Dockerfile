FROM golang:1-alpine AS builder

COPY . /go/src/github.com/mritd/notibot

WORKDIR /go/src/github.com/mritd/notibot

RUN set -ex \
    && apk add gcc musl-dev \
    && go install -trimpath -ldflags "-w -s"

FROM alpine

LABEL maintainer="mritd <mritd@linux.com>"

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN [ ! -e "/etc/nsswitch.conf" ] && echo 'hosts: files dns' > /etc/nsswitch.conf

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

ENTRYPOINT ["bash", "-c"]

CMD ["notibot --auth-mode ${NOTI_AUTH_MODE} --access-token ${NOTI_ACCESS_TOKEN} --username ${NOTI_USERNAME} --password ${NOTI_PASSWORD} --bot-api ${TELEGRAM_API} --bot-token ${TELEGRAM_TOKEN} --recipient ${TELEGRAM_RECIPIENT}"]