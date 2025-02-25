#!/bin/bash

set -u

# Resources:
# https://www.mongodb.com/docs/kafka-connector/current/source-connector/configuration-properties/
# https://www.mongodb.com/docs/kafka-connector/current/introduction/converters/#protobuf-converter
CONNECTOR_DATA=$(
	cat <<EOF
{
	"connector.class": "com.mongodb.kafka.connect.MongoSourceConnector",
	"tasks.max": "2",
	"connection.uri": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_CONNECTION_STRING",
	"database": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_DATABASE_NAME",
	"collection": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_MONGO_COLLECTION_NAME",
	"topic.prefix": "private.movie",
	"topic.suffix": "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_ID",
	"publish.full.document.only": true,
	"publish.full.document.only.tombstone.on.delete": false,
	"output.json.formatter": "com.mongodb.kafka.connect.source.json.formatter.SimplifiedJson",
	"output.format.key": "json",
	"key.converter.schemas.enable": false,
	"key.converter": "org.apache.kafka.connect.storage.StringConverter",
	"output.format.value": "schema",
	"output.schema.infer.value": true,
	"value.converter": "io.confluent.connect.protobuf.ProtobufConverter",
	"value.converter.schema.registry.url": "$SCHEMA_REGISTRY_URL"
}
EOF
)

EXACT_CONNECTOR_NAME=""$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_NAME"-"$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_ID""

# Run the migration
if [ "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_RUN_MIGRATION" == "true" ]; then
	echo "Running the migration for the movie service movies MongoDB source connector."
fi

# Delete the old connectors
if [ "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_DELETE_OLD_CONNECTORS" == "true" ]; then
	# If previous version of the connector (with the same prefix) is already registered, delete it
	res=$(curl -vS -X GET ""$CONNECT_URL"/connectors")
	if [ "$?" -ne 0 ]; then
		echo "Failed to get the list of connectors."
		exit 1
	fi
	previous_connector_names=$(echo "$res" | jq -c '.[]' | tr -d '"' | grep "$MOVIE_SERVICE_MOVIES_MONGODB_SOURCE_CONNECTOR_NAME" | grep -v "$EXACT_CONNECTOR_NAME")
	for connector_name in $previous_connector_names; do
		echo "Deleting the previous version of the connector: "$connector_name""
		curl -vS -X DELETE ""$CONNECT_URL"/connectors/"$connector_name""
		if [ "$?" -ne 0 ]; then
			echo "Failed to delete the previous version of the connector: "$connector_name""
			exit 1
		fi
	done
fi

# If the connector with exact name is not registered, register it, else update it
if [ $(curl -s -o /dev/null -w "%{http_code}" ""$CONNECT_URL"/connectors/$EXACT_CONNECTOR_NAME/status") -eq 404 ]; then
	CONNECTOR_DATA=$(
		cat <<EOF
{
    "name": "$EXACT_CONNECTOR_NAME",
    "config": $CONNECTOR_DATA
}
EOF
	)
	echo "Registering the movie service movies MongoDB source connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X POST ""$CONNECT_URL"/connectors" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to register the movie service movies MongoDB source connector."
		exit 1
	fi
	echo "Registered the movie service movies MongoDB source connector."
else
	echo "Updating the movie service movies MongoDB source connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X PUT ""$CONNECT_URL"/connectors/"$EXACT_CONNECTOR_NAME"/config" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to update the movie service movies MongoDB source connector."
		exit 1
	fi
	echo "Updated the movie service movies MongoDB source connector."
fi
