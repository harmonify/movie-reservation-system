# T001 - Initial Codebase Template

Build minimal codebase template on user service that will serve as a template for other upcoming services.

## Technical Requirements

- [x] Initial structure of `user-service`
- [x] Add `docker-compose.yml`
  - [x] Add dependencies
    - [x] PostgreSQL
    - [x] Redis
    - [x] Kafka
  - [x] Add observability
    - [x] Logging: Loki + Promtail + MinIO
    - [x] Tracing: Jaeger (Agent + Collector + Query) + Cassandra
    - [x] Metric: Prometheus
    - [x] Dashboard: Grafana
