# CQRS

Command Query Responsibility Segregation (CQRS) is a design pattern that segregates read and write operations for a data store into separate data models.

This allows each model to be optimized independently and can improve performance, scalability, and security of an application.

## Application

CQRS is applied for full-text search on movies and theaters.

## Storage

Search component uses Elasticsearch as its storage. It leverages Debezium as a CDC (change data capture) to pull in data from other data sources through Kafka, e.g. SQL database.

## References

- <https://microservices.io/patterns/data/cqrs.html>
- <https://microservices.io/patterns/data/transaction-log-tailing.html>
- <https://www.elastic.co/guide/en/elastic-stack/current/index.html>
- <https://debezium.io/documentation/reference/3.0/>
- <https://kafka.apache.org/documentation/>
