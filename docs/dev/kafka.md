# Kafka

## Golang SDK examples

Resource: <https://github.com/IBM/sarama/tree/main/examples>

## Topic name

General rules used to design topic name:

- Descriptive, general rule is the event's domain name followed by a past verb.
- Semantic version, in cases where the event schema changes.

The following is an example that satisfies both rules:

```txt
[domain name]-[past verb]-[semantic version]
```

Examples:

- `order_created_v1.0.0`, start.
- `order_paid_v1.0.0`, business flow succeed.
- `order_unpaid_v1.0.0`, business error.
- `order_failed_v1.0.0`, system errors.
- `order_completed_v1.0.0`, end.

## Topic ordering

For a ticket reservation process in a movie reservation system, maintaining correct event ordering is crucial. Kafka guarantees message order within a partition but not across partitions.

For ticket reservation process in this system, we can ensure ordering using movie `showtime_id`.

## Event message

Ensure that event messages have a consistent and extensible schema:

- Protobuf: For strong typing, schema evolution, and compact message sizes.
- Schema Registry: Confluent Schema Registry to enforce schema validation and versioning for event messages.
