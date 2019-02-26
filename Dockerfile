FROM alpine:3.9

WORKDIR /

COPY migrations /migrations

COPY StreamAnnotations-Backend /StreamAnnotations-Backend

ENTRYPOINT "./StreamAnnotations-Backend"