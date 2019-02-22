FROM alpine:3.9

WORKDIR /

COPY StreamAnnotations-Backend /StreamAnnotations-Backend

ENTRYPOINT "./StreamAnnotations-Backend"