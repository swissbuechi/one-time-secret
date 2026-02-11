FROM golang:1.26-alpine AS build
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3
ENV \
    VAULT_ADDR \
    VAULT_TOKEN \
    OTS_HTTP_REQUEST_LOG \
    OTS_HTTP_BINDING_ADDRESS \
    OTS_HTTPS_BINDING_ADDRESS \
    OTS_HTTPS_REDIRECT_ENABLED \
    OTS_TLS_AUTO_DOMAIN \
    OTS_TLS_CERT_FILEPATH \
    OTS_TLS_CERT_KEY_FILEPATH \
    OTS_VAULT_PREFIX

ENV TZ=Europe/Zurich
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /app
RUN mkdir -p /var/www
RUN apk add --no-cache ca-certificates curl libcap
RUN addgroup --gid 1000 -S app && adduser --uid 1000 -S app -G app
RUN chown -R app:app /app /var/www
COPY --from=build /app/app .
COPY ./static static
RUN setcap 'cap_net_bind_service=+ep' ./app
USER app
CMD [ "./app" ]
HEALTHCHECK --interval=5m --timeout=3s \
    CMD curl --fail --silent localhost${OTS_HTTP_BINDING_ADDRESS}/health | grep OK || exit 1