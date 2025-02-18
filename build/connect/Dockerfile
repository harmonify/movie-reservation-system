FROM confluentinc/cp-kafka-connect:7.8.1

# Confluent Hub: https://www.confluent.io/product/connectors/#oss-community-partner

# PostgreSQL source connector
RUN confluent-hub install --no-prompt debezium/debezium-connector-postgresql:2.5.4

# MongoDB source connector
RUN confluent-hub install --no-prompt mongodb/kafka-connect-mongodb:1.15.0

# Elasticsearch sink connector
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:14.1.2

# Run the Kafka Connect worker in distributed mode
CMD ["bash", "-c", "echo 'Starting Kafka Connect worker' && /etc/confluent/docker/run"]
