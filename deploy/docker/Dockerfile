# Stage 1: Runtime base with production system dependencies
FROM debian:12-slim AS runner

# Hardened minimal production environment for T-Bank Infrastructure
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libatomic1 \
    libstdc++6 \
    curl \
    procps \
    && rm -rf /var/lib/apt/lists/*

# Create unprivileged tbank user and service directories
RUN groupadd -g 10001 tbank && \
    useradd -u 10001 -g tbank -s /bin/false -d /opt/tbank tbank && \
    mkdir -p /opt/tbank/bin /opt/tbank/conf /opt/tbank/logs /opt/tbank/certs && \
    chown -R tbank:tbank /opt/tbank

WORKDIR /opt/tbank

# Copy SMT-verified super-optimized binaries and configuration
COPY --chown=tbank:tbank nginx.llm_opt.stripped /opt/tbank/bin/nginx
COPY --chown=tbank:tbank envoy.llm_opt.stripped /opt/tbank/bin/envoy
COPY --chown=tbank:tbank deploy/docker/entrypoint.sh /opt/tbank/bin/entrypoint.sh

RUN chmod +x /opt/tbank/bin/nginx /opt/tbank/bin/envoy /opt/tbank/bin/entrypoint.sh

# Production Environment variables for Zero-Overhead execution
ENV TBANK_SUPEROPT_VERIFIED="TRUE" \
    TBANK_GATEKEEPER_SMT="Z3_UNSAT" \
    TBANK_HUGEPAGES_ENABLED="1" \
    TBANK_IOURING_ENABLED="1"

EXPOSE 8080 8443 9090

USER tbank:tbank

ENTRYPOINT ["/opt/tbank/bin/entrypoint.sh"]
CMD ["nginx"]
