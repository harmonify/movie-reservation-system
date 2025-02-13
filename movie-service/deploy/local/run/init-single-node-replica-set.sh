#!/bin/bash

# Wait a bit to ensure that mongod is fully up
sleep 10

# Check if the replica set is already initiated
mongosh "mongodb://$MONGO_INITDB_ROOT_USERNAME:$MONGO_INITDB_ROOT_PASSWORD@movie-service-mongodb:27017" --quiet --eval "rs.status()" >/dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Replica set not initiated. Initiating now..."
	mongosh "mongodb://$MONGO_INITDB_ROOT_USERNAME:$MONGO_INITDB_ROOT_PASSWORD@movie-service-mongodb:27017" --quiet <<EOF
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "localhost:27017", priority: 1 },
  ]
});
EOF
	if [ $? -ne 0 ]; then
		echo "Failed to initiate replica set."
		exit 1
	fi
	echo "Replica set initiated."
else
	echo "Replica set already initiated."
fi
