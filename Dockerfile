FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY . /src/github.com/lileio/auth_service
ADD build/auth_service /bin
CMD ["auth_service", "server"]
