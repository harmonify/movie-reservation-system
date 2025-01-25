# Kafka

The following is the form of topic name used in this system:

```txt
[access level].[origin domain name].[past tense action verb].[schema major version]
```

## Components

### `[access level]`

This component specifies the access level of the topic, i.e. `public` and `private`. The access level is used to control who can access the topic. For example, `public` topics denotes that it is accessible to all services (e.g. `public.user.registered.v1`), while `private` topics denotes that it is only accessible to services within the same domain as the producer (e.g. `private.user.connect_configs.v1` used only for CDC configs that resides within the user domain).

### `[origin domain name]`

This identifies the system or domain where the data originates, such as "user", "movie", "order", etc.

### `[past tense action verb]`

This component specifies the business flow or event action, e.g. `registered`, `processed`, `paid`, etc.

### `[schema major version]`

This component specifies the schema version of the event message. The schema version is used to ensure backward compatibility when the schema changes. The schema version is incremented when the schema changes in a backward-incompatible way.

## DLQ topic naming

> Dead letter queue (or DLQ) serves as a holding area for messages that cannot be delivered or processed due to errors. By isolating these messages, a DLQ prevents them from disrupting the main queue's flow.

DLQ is ONLY used to hold messages that cannot be processed due to system errors (which differs from business flow error).

DLQ topics are named after their original topic, with the suffix `.dlq.[consumer group name]`. The consumer group name is used to identify the consumer group that consumes the original message but fails to process it.

## Test topic naming

Test topics are used for testing purposes and are named after their original topic, with the suffix `.test.[uuid]`. The UUID is used to ensure the tests are isolated from each other. These topics are created and deleted automatically by the test suite.

## Examples

The following are examples of topic naming in this system:

- `public.user.registered.v1`.
- `public.order.created.v1`, start.
- `public.order.paid.v1`, business flow.
- `public.order.unpaid.v1`, business flow error (not a system error).
- `public.order.completed.v1`, business flow.
- `public.order.processed.v1`, end.
- `private.order.created.v1.dlq.theater-service`, DLQ to hold system errors when theater service consumer group is processing the new order.
- `private.order.paid.v1.dlq.notification-service`, DLQ to hold system errors when notification service consumer group is processing the paid order.
- `private.order.created.v1.test.123e4567-e89b-12d3-a456-426614174000`, test topic for `order.created.v1`.

## Topic ordering

Maintaining correct event ordering is crucial for some cases, such as when dealing with finance transactions or reservations. Kafka guarantees message order within a partition but not across partitions.

### Ticket reservation process

To maintain correct ordering for the ticket reservation process in this system, we can use `showtime.id` to ensure we send messages to the same partition in `order.created.v1` topic.

## Event message

Ensure that event messages have a consistent and extensible schema:

- Protobuf: For strong typing, schema evolution, and compact message sizes.
- Schema Registry: Confluent Schema Registry to enforce schema validation and versioning for event messages.

## Resources

- [Kafka topic naming convention](https://www.confluent.io/learn/kafka-topic-naming-convention/)
- [Kafka Go client examples](https://github.com/IBM/sarama/tree/main/examples)
- [Building resilient microservices using Kafka DLQ](https://medium.com/@Games24x7Tech/building-resilient-micro-services-using-kafka-dlq-5654faee6de2)
