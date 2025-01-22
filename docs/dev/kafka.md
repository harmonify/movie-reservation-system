# Kafka

The following is the form of topic name used in this system:

```txt
[origin domain name].[past tense action verb].[schema major version]
```

## Components

### `[origin domain name]`

This identifies the system or domain where the data originates, such as "user", "movie", "order", etc.

### `[past tense action verb]`

This component specifies the business flow or event action, i.e. "created", "processed", "paid", etc.

### `[schema major version]`

This component specifies the schema version of the event message. The schema version is used to ensure backward compatibility when the schema changes. The schema version is incremented when the schema changes in a backward-incompatible way.

## DLQ topic naming

> Dead letter queue (or DLQ) serves as a holding area for messages that cannot be delivered or processed due to errors. By isolating these messages, a DLQ prevents them from disrupting the main queue's flow.

DLQ is ONLY used to hold messages that cannot be processed due to system errors (which differs from business flow error).

DLQ topics are named after their original topic, with the suffix `.dlq`.

## Test topic naming

Test topics are used for testing purposes and are named after their original topic, with the suffix `.test.[uuid]`. The UUID is used to ensure the tests are isolated from each other. These topics are created and deleted automatically by the test suite.

## Command topic naming

Command topics are used to send commands to services. Although they are somewhat contradictory to the event-driven architecture, they are still useful in some cases, i.e. where we need to not over

They are named differently from event topics. The naming convention for command topics is:

```txt
[domain name].[command name].[schema major version]
```

## Examples

The following are examples of topic naming in this system:

- `order.created.v1`, start.
- `order.paid.v1`, business flow.
- `order.unpaid.v1`, business flow error (not a system error).
- `order.completed.v1`, business flow.
- `order.processed.v1`, end.
- `order.created.v1.dlq`, DLQ to hold system errors when processing new order.
- `order.paid.v1.dlq`, DLQ to hold system errors when processing paid order.
- `order.created.v1.test.123e4567-e89b-12d3-a456-426614174000`, test topic for `order.created.v1`.

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
