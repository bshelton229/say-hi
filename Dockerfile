FROM alpine

RUN apk add --update \
        jq \
        vim \
        curl \
        py-pip \
    && pip install \
        awscli \
    && rm -rf /var/cache/apk/*

COPY ./say-hi-linux /usr/local/bin/say-hi

EXPOSE 8082
CMD ["/usr/local/bin/say-hi"]
