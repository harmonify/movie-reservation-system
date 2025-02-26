#!/bin/bash

set -u

TABLE_NAME=""$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_SCHEMA"."$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_TABLE""

OUTBOX_TABLE_MIGRATION=$(
	cat <<EOF
START TRANSACTION;

CREATE TABLE IF NOT EXISTS "$TABLE_NAME" (
	id UUID PRIMARY KEY NOT NULL,
	aggregatetype VARCHAR(255) NOT NULL,
    aggregateid VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,
    tracingspancontext JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "$TABLE_NAME".id IS 'The tracing ID or request ID associated with the event';

COMMENT ON COLUMN "$TABLE_NAME".aggregatetype IS 'The aggregate event type, e.g., registered';

COMMENT ON COLUMN "$TABLE_NAME".payload IS 'The outbox payload containing the event data in Protobuf binary format';

COMMENT ON COLUMN "$TABLE_NAME".tracingspancontext IS 'The tracing span context associated with the event';

COMMIT;

EOF
)

# Resources:
# https://debezium.io/documentation/reference/3.0/connectors/postgresql.html
# https://debezium.io/documentation/reference/3.0/transformations/outbox-event-router.html#basic-outbox-table
CONNECTOR_DATA=$(
	cat <<EOF
{
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_HOST",
    "database.port": "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_PORT",
    "database.user": "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_USER",
    "database.password": "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_PASSWORD",
    "database.dbname": "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_NAME",
    "table.field.event.id": "id",
    "table.field.event.key": "aggregateid",
    "table.field.event.payload": "payload",
    "table.include.list": "$TABLE_NAME",
    "tasks.max": "2",
    "tombstones.on.delete": "false",
    "topic.prefix": "private",
    "tracing.operation.name": "debezium-read",
    "tracing.span.context.field": "context",
    "tracing.with.context.field.only": "true",
    "transforms": "outbox",
    "transforms.outbox.route.topic.replacement": "public.user.\${routedByValue}.$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_CONNECTOR_ID",
    "transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
    "route.by.field": "aggregatetype",
    "value.converter": "io.debezium.converters.BinaryDataConverter"
}
EOF
)

EXACT_CONNECTOR_NAME=""$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_CONNECTOR_NAME"-"$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_CONNECTOR_ID""

# Run the migration
if [ "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_RUN_MIGRATION" == "true" ]; then
	# Create pass file
	PGPASSFILE=$(mktemp)
	echo "*:*:*:"$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_USER":"$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_PASSWORD >$PGPASSFILE
	export PGPASSFILE
	# Create the outbox table
	PGPASSWORD="$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_PASSWORD" psql -h "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_HOST" \
		-p "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_PORT" \
		-U "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_USER" \
		-d "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DB_NAME" \
		-c "$OUTBOX_TABLE_MIGRATION"
	if [ "$?" -ne 0 ]; then
		echo "Failed to create the outbox table."
		rm $PGPASSFILE
		exit 1
	fi
	rm $PGPASSFILE
fi

# Delete the old connectors
if [ "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_DELETE_OLD_CONNECTORS" == "true" ]; then
	# If previous version of the connector (with the same prefix) is already registered, delete it
	res=$(curl -vS -X GET ""$CONNECT_URL"/connectors")
	if [ "$?" -ne 0 ]; then
		echo "Failed to get the list of connectors."
		exit 1
	fi
	echo "Previous connectors: "$res""
	previous_connector_names=$(echo "$res" | jq -c '.[]' | tr -d '"' | grep "$USER_SERVICE_OUTBOX_POSTGRESQL_SOURCE_CONNECTOR_NAME" | grep -v "$EXACT_CONNECTOR_NAME")
	for connector_name in $previous_connector_names; do
		echo "Deleting the previous version of the connector: "$connector_name""
		curl -vS -X DELETE ""$CONNECT_URL"/connectors/"$connector_name""
		if [ "$?" -ne 0 ]; then
			echo "Failed to delete the previous version of the connector: "$connector_name""
			exit 1
		fi
	done
fi

# Check if the connector is already registered
if [ $(curl -s -o /dev/null -w "%{http_code}" ""$CONNECT_URL"/connectors/"$EXACT_CONNECTOR_NAME"/status") -eq 404 ]; then
	CONNECTOR_DATA=$(
		cat <<EOF
{
    "name": "$EXACT_CONNECTOR_NAME",
    "config": $CONNECTOR_DATA
}
EOF
	)
	echo "Registering the user service outbox PostgreSQL source connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X POST ""$CONNECT_URL"/connectors/" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to register the user service outbox PostgreSQL source connector."
		exit 1
	fi
	echo "Registered the user service outbox PostgreSQL source connector."
else
	echo "Updating the user service outbox PostgreSQL source connector with the following configuration:"
	echo "$CONNECTOR_DATA"
	curl -vS -X PUT ""$CONNECT_URL"/connectors/"$EXACT_CONNECTOR_NAME"/config" -H "Content-Type: application/json" -d "$CONNECTOR_DATA"
	if [ "$?" -ne 0 ]; then
		echo "Failed to update the user service outbox PostgreSQL source connector."
		exit 1
	fi
	echo "Updated the user service outbox PostgreSQL source connector."
fi

echo
