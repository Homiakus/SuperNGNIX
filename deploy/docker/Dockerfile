# Stage 1: Runtime base with production system dependencies
FROM debian:12-slim AS runner

# Hardened minimal production environment for EdgeGateway Infrastructure
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libatomic1 \
    libstdc++6 \
    curl \
    procps \
    && rm -rf /var/lib/apt/lists/*

# Create unprivileged edgegate user and service directories
RUN groupadd -g 10001 edgegate && \
    useradd -u 10001 -g edgegate -s /bin/false -d /opt/edgegate edgegate && \
    mkdir -p /opt/edgegate/bin /opt/edgegate/conf /opt/edgegate/logs /opt/edgegate/certs && \
    chown -R edgegate:edgegate /opt/edgegate

WORKDIR /opt/edgegate

# Copy SMT-verified super-optimized binaries and configuration
COPY --chown=edgegate:edgegate nginx.llm_opt.stripped /opt/edgegate/bin/nginx
COPY --chown=edgegate:edgegate envoy.llm_opt.stripped /opt/edgegate/bin/envoy
COPY --chown=edgegate:edgegate deploy/docker/entrypoint.sh /opt/edgegate/bin/entrypoint.sh

RUN chmod +x /opt/edgegate/bin/nginx /opt/edgegate/bin/envoy /opt/edgegate/bin/entrypoint.sh

# Production Environment variables for Zero-Overhead execution
ENV EDGEGATE_SUPEROPT_VERIFIED="TRUE" \
    EDGEGATE_GATEKEEPER_SMT="Z3_UNSAT" \
    EDGEGATE_HUGEPAGES_ENABLED="1" \
    EDGEGATE_IOURING_ENABLED="1"

EXPOSE 8080 8443 9090

USER edgegate:edgegate

ENTRYPOINT ["/opt/edgegate/bin/entrypoint.sh"]
CMD ["nginx"]
