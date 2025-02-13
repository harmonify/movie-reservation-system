#!/bin/bash

# Wait a bit to ensure that mongod is fully up
sleep 10

# Check if the replica set is already initiated
mongosh "mongodb://$MONGO_INITDB_ROOT_USERNAME:$MONGO_INITDB_ROOT_PASSWORD@localhost:27017" --quiet --eval "rs.status()" > /dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "Replica set not initiated. Initiating now..."
  mongosh  "mongodb://$MONGO_INITDB_ROOT_USERNAME:$MONGO_INITDB_ROOT_PASSWORD@localhost:27017" --quiet <<EOF
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "host.docker.internal:27017", priority: 1 },
    { _id: 1, host: "host.docker.internal:27018", priority: 0.5 },
    { _id: 2, host: "host.docker.internal:27019", arbiterOnly: true }
  ]
});
EOF
  echo "Replica set initiated."
else
  echo "Replica set already initiated."
fi

tail -f /dev/null
