FROM alpine:3.10

WORKDIR /app

RUN apk --no-cache add tzdata bash ca-certificates mailcap wkhtmltopdf \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/*

COPY ./build/docker/app/files /
COPY ./.artifacts .

RUN chmod +x /docker/bin/* \
    && chmod +x bin

ENTRYPOINT ["/docker/bin/entrypoint.sh"]
