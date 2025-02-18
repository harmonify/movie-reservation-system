#!/bin/bash

set -u

SOURCE_TOPIC="private.movie."$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_DATABASE_NAME"."$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_COLLECTION_NAME".$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_ID"

# Resources:
# https://docs.confluent.io/kafka-connectors/elasticsearch/14.1/overview.html
#
# Misc:
# "value.converter.basic.auth.credentials.source": "$BASIC_AUTH_CREDENTIALS_SOURCE",
# "value.converter.basic.auth.user.info": "$SCHEMA_REGISTRY_BASIC_AUTH_USER_INFO",
#
# The transformation chain:
# 1. Rename the field "id" to "_id" in the value.
# 2. Convert the value to the key using the field "id".
# 3. Extract the field "id" from the key.
CONNECTOR_DATA=$(
	cat <<EOF
{
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "tasks.max": "2",
    "topics": "$SOURCE_TOPIC",
    "connection.url": "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_CONNECTION_URL",
    "connection.username": "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_USERNAME",
    "connection.password": "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_PASSWORD",
	"key.converter":"org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
	"value.converter.schema.registry.url": "$SCHEMA_REGISTRY_URL",
    "transforms": "renameIdField,copyIdFieldToKeyDoc,extractActualKeyFromIdFieldInKeyDoc,removeIdField",
    "transforms.renameIdField.type": "org.apache.kafka.connect.transforms.ReplaceField\$Value",
    "transforms.renameIdField.renames": "_id:id",
    "transforms.copyIdFieldToKeyDoc.type":"org.apache.kafka.connect.transforms.ValueToKey",
    "transforms.copyIdFieldToKeyDoc.fields":"id",
    "transforms.extractActualKeyFromIdFieldInKeyDoc.type":"org.apache.kafka.connect.transforms.ExtractField\$Key",
    "transforms.extractActualKeyFromIdFieldInKeyDoc.field":"id",
    "transforms.removeIdField.type": "org.apache.kafka.connect.transforms.ReplaceField\$Value",
    "transforms.removeIdField.blacklist": "id",
    "write.method":"upsert",
	"read.timeout.ms": "8000",
	"batch.size": "1000",
	"max.buffered.records": "10000"
}
EOF
)

EXACT_CONNECTOR_NAME=""$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_CONNECTOR_NAME"-"$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_CONNECTOR_ID""

# Run the migration
if [ "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_RUN_MIGRATION" == "true" ]; then
	echo "Running the migration for the movie search service movies Elasticsearch sink connector."
fi

# Delete the old connectors
if [ "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_DELETE_OLD_CONNECTORS" == "true" ]; then
	# If previous version of the connector (with the same prefix) is already registered, delete it
	res=$(curl -vS -X GET ""$CONNECT_URL"/connectors")
	if [ "$?" -ne 0 ]; then
		echo "Failed to get the list of connectors."
		exit 1
	fi

	previous_connector_names=$(echo "$res" | jq -c '.[]' | tr -d '"' | grep "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_CONNECTOR_NAME" | grep -v "$EXACT_CONNECTOR_NAME")
	for connector_name in $previous_connector_names; do
		echo "Deleting the previous version of the connector: "$connector_name""
		curl -vS -X DELETE ""$CONNECT_URL"/connectors/"$connector_name""
		if [ "$?" -ne 0 ]; then
			echo "Failed to delete the previous version of the connector: "$connector_name""
			exit 1
		fi
	done
fi

# Delete the old index
if [ "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_DELETE_OLD_INDEX" == "true" ]; then
	echo "Deleting the old index for the movie search service movies Elasticsearch sink connector."

	# If previous version of the index (with the same prefix) is already registered, delete it
	MATCHING_INDEX_NAME=$(curl -sS -X GET ""$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_CONNECTION_URL"/_cat/indices" -u "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_USERNAME":"$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_PASSWORD" | awk '{print $3}' | grep "$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_CONNECTOR_NAME" | grep -v "$EXACT_CONNECTOR_NAME")
	if [ ! -z "$MATCHING_INDEX_NAME" ]; then
		echo "Deleting the previous version of the index: "$MATCHING_INDEX_NAME""
		curl -vS -X DELETE ""$MOVIE_SEARCH_SERVICE_MOVIES_ELASTICSEARCH_SINK_ELASTICSEARCH_CONNECTION_URL"/"$MATCHING_INDEX_NAME""
		if [ "$?" -ne 0 ]; then
			echo "Failed to delete the previous version of the index: "$MATCHING_INDEX_NAME""
			exit 1
		fi
	fi
fi

# If the connector with exact name is not registered, register it, else update it
if [ $(curl -s -o /dev/null -w "%{http_code}" ""$CONNECT_URL"/connectors/"$EXACT_CONNECTOR_NAME"/status") -eq 404 ]; then
	# Register the connector
	CONNECTOR_DATA=$(
		cat <<EOF
{
    "name": "$EXACT_CONNECTOR_NAME",
    "config": $CONNECTOR_DATA
}
EOF
	)
	echo "Registering the movie service movies Elasticsearch sink connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X POST ""$CONNECT_URL"/connectors" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to register the movie service movies Elasticsearch sink connector."
		exit 1
	fi
	echo "Registered the movie service movies Elasticsearch sink connector."
else
	echo "Updating the movie service movies Elasticsearch sink connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X PUT ""$CONNECT_URL"/connectors/"$EXACT_CONNECTOR_NAME"/config" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to update the movie service movies Elasticsearch sink connector."
		exit 1
	fi
	echo "Updated the movie service movies Elasticsearch sink connector."
fi
