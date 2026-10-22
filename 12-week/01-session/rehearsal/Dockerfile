# The artifact of the rehearsal: built ONCE per version, then promoted by tag. The environment
# (ENVIRONMENT) is read when the container starts, never at build time.
FROM nginx:1.27-alpine
ARG VERSION
LABEL org.opencontainers.image.version="${VERSION}"
ENV APP_VERSION="${VERSION}"
COPY app/default.conf.template /etc/nginx/templates/default.conf.template
HEALTHCHECK --interval=2s --timeout=2s --start-period=2s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1
