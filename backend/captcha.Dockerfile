# multi-stage builds

# step 1
FROM golang:1.19.3-alpine3.16 as builder

# build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN mkdir -p /app/src/webtoy_base
RUN mkdir -p /app/src/webtoy_msg_captcha
RUN mkdir -p /app/src/webtoy_captcha
RUN mkdir -p /app/bin
COPY ./webtoy_base/ /app/src/webtoy_base/
COPY ./webtoy_msg_captcha/ /app/src/webtoy_msg_captcha
COPY ./webtoy_captcha/ /app/src/webtoy_captcha/
WORKDIR /app/src/webtoy_captcha
RUN go build -o /app/bin/webtoy_captcha

# step 2
FROM alpine:3.16

# run
RUN mkdir -p /app/bin
COPY --from=builder /app/bin/ /app/bin
EXPOSE 8080
ENV PATH="/app/bin:${PATH}"
WORKDIR /app/bin
CMD webtoy_captcha
