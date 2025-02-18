#!/bin/bash

set -euo pipefail

# Resources:
# https://www.mongodb.com/docs/kafka-connector/current/source-connector/configuration-properties/
# https://www.mongodb.com/docs/kafka-connector/current/introduction/converters/#protobuf-converter
CONNECTOR_DATA=$(
	cat <<EOF
{
	"connector.class": "com.mongodb.kafka.connect.MongoSourceConnector",
	"connection.uri": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_CONNECTION_STRING",
	"database": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_DATABASE_NAME",
	"collection": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_COLLECTION_NAME",
	"topic.prefix": "private.movie",
	"topic.suffix": "v1",
	"publish.full.document.only": true,
	"publish.full.document.only.tombstone.on.delete": false,
	"output.format.key": "json",
	"key.converter.schemas.enable": false,
	"key.converter": "org.apache.kafka.connect.storage.StringConverter",
	"output.format.value": "schema",
	"output.schema.infer.value": true,
	"value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
	"value.converter.schema.registry.url": "$SCHEMA_REGISTRY_URL",
	"tasks.max": "2"
}
EOF
)

echo "Registering the movie service movies MongoDB source connector with the following configuration:"
echo "$CONNECTOR_DATA"

# Check if the connector is already registered
if [ $(curl -s -o /dev/null -w "%{http_code}" ""$CONNECT_URL"/connectors/$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_NAME/status") -eq 404 ]; then
	# Register the connector
	CONNECTOR_DATA=$(
		cat <<EOF
{
    "name": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_NAME",
    "config": $CONNECTOR_DATA
}
EOF
	)
	curl -vS -X POST ""$CONNECT_URL"/connectors" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
else
	# Update the connector
	curl -vS -X PUT ""$CONNECT_URL"/connectors/$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_NAME/config" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
fi

echo
echo "Done registering the movie service movies MongoDB source connector."
