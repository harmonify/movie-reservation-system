#!/bin/bash

set -euo pipefail

TOPIC="private.movie."$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_DATABASE_NAME"."$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_COLLECTION_NAME".v1"

# "value.converter.basic.auth.credentials.source": "$BASIC_AUTH_CREDENTIALS_SOURCE",
# "value.converter.basic.auth.user.info": "$SCHEMA_REGISTRY_BASIC_AUTH_USER_INFO",
CONNECTOR_DATA=$(
	cat <<EOF
{
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "topics": "$TOPIC",
    "connection.url": "$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_ELASTICSEARCH_CONNECTION_URL",
    "connection.username": "$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_ELASTICSEARCH_USERNAME",
    "connection.password": "$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_ELASTICSEARCH_PASSWORD",
    "tasks.max": "2",
	"key.converter":"org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
	"value.converter.schema.registry.url": "$SCHEMA_REGISTRY_URL",
    "transforms": "renameIdField,valueToKey,extractIdField",
    "transforms.renameIdField.type": "org.apache.kafka.connect.transforms.ReplaceField\$Value",
    "transforms.renameIdField.renames": "_id:id",
    "transforms.valueToKey.type":"org.apache.kafka.connect.transforms.ValueToKey",
    "transforms.valueToKey.fields":"id",
    "transforms.extractIdField.type":"org.apache.kafka.connect.transforms.ExtractField\$Key",
    "transforms.extractIdField.field":"id",
    "write.method":"upsert"
}
EOF
)

echo "Registering the movie service movies Elasticsearch sink connector with the following configuration:"
echo "$CONNECTOR_DATA"

# Check if the connector is already registered
if [ $(curl -s -o /dev/null -w "%{http_code}" ""$CONNECT_URL"/connectors/"$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_CONNECTOR_NAME"/status") -eq 404 ]; then
	# Register the connector
	CONNECTOR_DATA=$(
		cat <<EOF
{
    "name": "$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_CONNECTOR_NAME",
    "config": $CONNECTOR_DATA
}
EOF
	)
	curl -vS -X POST ""$CONNECT_URL"/connectors" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
else
	# Update the connector
	curl -vS -X PUT ""$CONNECT_URL"/connectors/"$MOVIE_SEARCH_SERVICE_ELASTICSEARCH_SINK_CONNECTOR_NAME"/config" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
fi

echo
echo "Done registering the movie service movies Elasticsearch sink connector."
