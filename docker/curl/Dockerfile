FROM alpine:edge

RUN apk update && apk add --no-cache \
    curl

CMD ["sh", "-c", "curl ${PARAMS} ${HOST}"]
