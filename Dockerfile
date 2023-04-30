FROM golang:1.20-alpine AS build
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3
ENV \
    VAULT_ADDR \
    VAULT_TOKEN \
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
RUN apk add --no-cache ca-certificates curl libcap
RUN addgroup --gid 1000 -S app && adduser --uid 1000 -S app -G app
RUN chown -R app:app /app
COPY --from=build /app/app .
ADD ./static static
RUN setcap 'cap_net_bind_service=+ep' ./app
USER app
CMD [ "./app" ]
# HEALTHCHECK CMD curl --fail --silent localhost:80/health | grep UP || exit 1