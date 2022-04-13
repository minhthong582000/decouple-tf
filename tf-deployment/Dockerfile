# Builder
FROM golang:1.17.8-alpine3.15 as builder

RUN apk update && apk upgrade && \
    apk --update add git make npm curl

RUN npm install --global cdktf-cli@latest

ARG terraform_version=1.1.6
RUN cd /tmp \
    && curl -LO https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip \
    && unzip terraform_${terraform_version}_linux_amd64.zip \
    && mv terraform /usr/local/bin \
    && rm terraform_${terraform_version}_linux_amd64.zip

WORKDIR /app

COPY . .

RUN cdktf get
RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata nodejs && \
    mkdir /app 

WORKDIR /app 

COPY --from=builder /app/engine /app
COPY --from=builder /usr/local/bin/terraform /usr/local/bin/

CMD /app/engine