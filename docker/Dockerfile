FROM alpine:3.20

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="GeekWorkCode <asd486@live.com>" \
  org.label-schema.name="Plane API CLI" \
  org.label-schema.vendor="GeekWorkCode" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/GeekWorkCode/go-plane
LABEL org.opencontainers.image.description="Plane API CLI"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/go-plane /bin/

ENTRYPOINT ["/bin/go-plane"] 