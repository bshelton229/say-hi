FROM alpine

RUN apk add --update \
        jq \
        vim \
        curl \
        py-pip \
        bash \
        openssl \
    && pip install \
        awscli \
    && rm -rf /var/cache/apk/*

# Install kubectl and helm
RUN cd /usr/local/bin \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod 755 /usr/local/bin/kubectl \
    && cd /tmp \
    && curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh \
    && chmod 755 get_helm.sh \
    && ./get_helm.sh \
    && rm ./get_helm.sh

COPY ./say-hi-linux /usr/local/bin/say-hi

EXPOSE 8082
CMD ["/usr/local/bin/say-hi"]
