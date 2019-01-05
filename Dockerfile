FROM alpine:3.8

WORKDIR /

ENV MYSQL_USER=root
ENV MYSQL_PASS=""
ENV MYSQL_ADDR="localhost:3306"
ENV MYSQL_DBNAME="infoboxes"

COPY InfoBoxes-Backend /InfoBoxes-Backend

ENTRYPOINT "./InfoBoxes-Backend"