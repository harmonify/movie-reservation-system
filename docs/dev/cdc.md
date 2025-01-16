# CDC

## Debezium

### Set up PostgreSQL source connector

#### Execute `docker compose up -d`

#### Confirm PostgreSQL database configuration

In the PostgreSQL client, execute following queries:

```sql
show wal_level; -- logical
show max_replication_slots; -- 10
show max_wal_senders; -- 10
show shared_preload_libraries; -- decoderbufs
```

#### Perform database migrations

```sh
cd user-service
make migration:postgresql:up
```

#### Register PostgreSQL source connector

```sh
cd user-service
make debezium:register-postgresql-source-connector
```

#### Perform INSERT/UPDATE/DELETE operation within the PostgreSQL database

#### Visit Kafdrop UI

Visit <http://localhost:9001/topic/user_connect.public.users/allmessages> to view all topic messages produced by Debezium PostgreSQL source connector.

Here is [an example](./assets/data/cdc-debezium-message-example.json) of such message.

## Resources

- [Official Debezium PostgreSQL source connector documentation](https://debezium.io/documentation/reference/stable/connectors/postgresql.html)
- [Tutorial Debezium PostgreSQL with decoderbufs (default plugin)](https://dev.to/emtiajium/track-every-postgresql-data-change-using-debezium-5e19)
- [Tutorial Debezium PostgreSQL with pgoutput](https://medium.com/@arijit.mazumdar/beyond-the-basics-of-debezium-for-postgresql-part-1-d1c6952ae110)
- [Alternative (Confluent Platform AIO)](https://github.com/confluentinc/cp-all-in-one/blob/latest/cp-all-in-one/docker-compose.yml)
- [Tutorial Confluent Elasticsearch sink connector](https://www.confluent.io/blog/kafka-elasticsearch-connector-tutorial/)
- [Tutorial Debezium PostgreSQL source connector + Confluent Elasticsearch sink connector](https://medium.com/@cagataygokcel/real-time-data-streaming-from-postgresql-to-elasticsearch-via-kafka-and-debezium-b624b43cadb)
