# From alpine
# Install curl
# Install PostgreSQL client
# Install MongoDB client

FROM alpine:3

ARG CONNECTOR_REGISTRAR_VERSION=1.0.0

RUN apk add --no-cache bash curl jq
RUN apk add --no-cache postgresql-client
RUN apk add --no-cache mongodb-tools

# Default runtime environment variables
ENV CONNECTOR_REGISTRAR_VERSION=${CONNECTOR_REGISTRAR_VERSION}

RUN mkdir /connector-registrar

WORKDIR /connector-registrar

ADD ../*.sh .

RUN find . -type f -name "*.sh" -exec chmod +x {} \;

CMD [ "bash", "/connector-registrar/run.sh" ]
