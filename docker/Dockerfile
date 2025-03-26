FROM alpine:3.20

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Your Name <your.email@example.com>" \
  org.label-schema.name="Plane API CLI" \
  org.label-schema.vendor="Your Name" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/yourname/go-plane
LABEL org.opencontainers.image.description="Plane API CLI"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/go-plane /bin/

ENTRYPOINT ["/bin/go-plane"] 