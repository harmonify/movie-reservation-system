# T001 - Initial Codebase Template

Build minimal codebase template on user service that will serve as a template for other upcoming services.

## Technical Requirements

- [x] Initial structure of `user-service`
- [x] Add infrastructure code on `docker-compose.yml`
  - [x] OLTPs
    - [x] PostgreSQL
    - [x] Redis
    - [x] MongoDB
    - [x] MySQL
  - [ ] OLAPs
    - [ ] ClickHouse
    - [x] Elasticsearch
  - [x] Inter-service communication
    - [x] gRPC
    - [x] Kafka
  - [x] ETLs
    - [x] Apache Kafka Connect (its connectors must be OSS)
  - [x] Monitoring w/ OpenTelemetry
    - [x] Logging: Loki + Promtail + MinIO
    - [x] Tracing: Jaeger (Agent + Collector + Query) + Cassandra
    - [x] Metric: Prometheus
    - [x] Dashboard: Grafana
    - [ ] Alerting
    - [ ] Health Check
  - [x] Security
    - [x] OAuth2
    - [x] JWT
    - [x] HTTPS
    - [x] API Gateway
  - [ ] CI/CD Pipeline w/ Github Actions
    - [ ] Code Quality
      - [ ] SonarQube
    - [ ] Code Review
      - [ ] Pull Request Template
    - [ ] Continuous Deployment
      - [ ] Docker Registry
      - [ ] Helm for Kubernetes
      - [x] ~~Production deployment~~
  - [ ] Documentation
    - [ ] Swagger
    - [x] PlantUML
  - [ ] API Gateway
    - [ ] Traefik
  - [ ] Service Mesh (& Service Discovery)
    - [ ] Istio / Linkerd / Consul
  - [ ] Container Orchestration
    - [ ] Kubernetes
  - [ ] Configuration Management
    - [ ] Consul
  - [ ] Secret Management
    - [ ] Vault
