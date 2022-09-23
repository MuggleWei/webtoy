FROM golang:1.19.1-alpine

# build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN mkdir -p /app/src/webtoy_base
RUN mkdir -p /app/src/webtoy_captcha
RUN mkdir -p /app/bin
COPY ./webtoy_base/ /app/src/webtoy_base/
COPY ./webtoy_captcha/ /app/src/webtoy_captcha/
WORKDIR /app/src/webtoy_captcha
RUN go build -o /app/bin/webtoy_captcha

# clean source
RUN rm -rf /app/src

# run
EXPOSE 8080
ENV PATH="/app/bin:${PATH}"
WORKDIR /app/bin
CMD webtoy_captcha
