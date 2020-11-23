FROM golang:alpine as builder



RUN apk add --no-cache curl tzdata

WORKDIR /app/svc
COPY ./dist/playground-linux-x64 .

CMD ./playground-linux-x64
