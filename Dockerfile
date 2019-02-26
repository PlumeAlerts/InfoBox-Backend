FROM alpine:3.9

RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /

COPY migrations /migrations

COPY StreamAnnotations-Backend /StreamAnnotations-Backend

ENTRYPOINT "./StreamAnnotations-Backend"