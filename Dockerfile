FROM alpine

RUN apk add --update \
        jq \
        vim \
        curl \
        py-pip \
        bash \
    && pip install \
        awscli \
    && rm -rf /var/cache/apk/*

# Install kubectl for testing
RUN cd /usr/local/bin \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod 755 /usr/local/bin/kubectl

COPY ./say-hi-linux /usr/local/bin/say-hi

EXPOSE 8082
CMD ["/usr/local/bin/say-hi"]
