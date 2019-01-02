FROM golang:1.11.2-alpine

RUN apk add --update netcat-openbsd bash && rm -rf /var/cache/apk/*

WORKDIR /

ENV GO111MODULE=on

ENV MYSQL_USER=root
ENV MYSQL_PASS=""
ENV MYSQL_ADDR="localhost:3306"
ENV MYSQL_DBNAME="infoboxes"

COPY InfoBoxes-Backend /InfoBoxes-Backend

ENTRYPOINT "./InfoBoxes-Backend"