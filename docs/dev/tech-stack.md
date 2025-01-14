# Tech Stack (WIP)

## Infrastructure

- Kubernetes, for:
  - load balancer
  - service discovery
- Kafka, for message broker
  - Protobuf, for event message schema
  - Confluent Control Center, for schema registry

## Microservices

### user-service

- Go
- GORM
- Casbin
- PostgreSQL
- Redis

### movie-service

- Go
- MongoDB
- [Debezium](https://debezium.io/documentation/reference/stable/install.html) (MongoDB -> movie-query-service's Elasticsearch)
- Redis

### movie-query-service

- Go
- Elasticsearch

### theater-service

- Go
- GORM
- PostgreSQL
- Elasticsearch ?
- Redis

### ticket-service

- Go
- GORM
- PostgreSQL
- Redis

### order-service

- Go
- GORM
- PostgreSQL
- Debezium (PostgreSQL -> order-query-service's Elasticsearch)
- Redis

### order-query-service

- Go
- Elasticsearch

## FAAS

### backup-user-service-postgresql

- Go

### backup-movie-service-mongodb

- Go

### backup-theater-service-postgresql

- Go

### backup-ticket-service-postgresql

- Go

### backup-order-service-postgresql

- Go

### process-new-order (orchestrator job)

- Go

## CLI

Help to control and maintain system.

## Client

- Typescript
- React.js
- Next.js
- Tailwind

## Notification

### Email

- Mailgun

### SMS

- Twilio

## Analytic

### User activity analytic

- Posthog

## Observability

- Grafana

### Log Aggregation

- Loki
- Promtail
- Min.IO, for storage

### Distributed Tracing

- Jaeger
- Cassandra

### Application Metric

- Prometheus
- PostgreSQL exporter
- Redis exporter
- Elasticsearch exporter
- MongoDB exporter
- Kafka exporter

### Alerting

For developers.

#### Slack

This will include link to Grafana, which includes error log and trace for debugging.
