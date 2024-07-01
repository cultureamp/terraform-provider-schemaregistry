# Base stage
FROM centos:7 as base

USER root

# Install necessary tools
RUN yum install -y yum-utils \
    && yum install -y make wget \
    && yum clean all \
    && rm -rf /var/cache/yum

USER appuser

# Final stage
FROM confluentinc/cp-schema-registry:7.6.1

USER root

# Copy tools from the base stage
COPY --from=base /usr/bin/make /usr/bin/make
COPY --from=base /usr/bin/wget /usr/bin/wget

# Run yum update to patch any security vulnerabilities, excluding the Confluent repository
RUN yum install -y yum-utils \
    && yum-config-manager --save --setopt=ubi-8-baseos-rpms.sslverify=false \
    && yum update --disablerepo=Confluent -y \
    && yum clean all \
    && rm -rf /var/cache/yum

USER appuser
