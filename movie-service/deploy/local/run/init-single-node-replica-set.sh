#!/bin/bash

set -euo pipefail

if [ -z "$MONGO_URI" ]; then
	echo "MONGO_URI is not set."
	exit 1
fi

if [ -z "$MONGO_REPLICA_SET_MEMBER" ]; then
	MONGO_REPLICA_SET_MEMBER="host.docker.internal:27017"
fi

# Wait a bit to ensure that mongod is fully up
sleep 2

# Check if the replica set is already initiated
mongosh "$MONGO_URI" --quiet --eval "rs.status()" >/dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Replica set not initiated. Initiating now..."
	mongosh "$MONGO_URI" --quiet <<EOF
rs.initiate({
  _id: "$MONGO_REPLICA_SET_NAME",
  members: [
    { _id: 0, host: "$MONGO_REPLICA_SET_MEMBER", priority: 1 },
  ]
});
EOF
	if [ $? -ne 0 ]; then
		echo "Failed to initiate replica set."
		exit 1
	fi
	echo "Replica set initiated."
else
	# update the replica set configuration
	echo "Replica set already initiated. Updating configuration..."
	mongosh "$MONGO_URI" --quiet <<EOF
cfg = rs.conf();
cfg.members[0].host = "$MONGO_REPLICA_SET_MEMBER";
rs.reconfig(cfg, { force: true });
EOF
	if [ $? -ne 0 ]; then
		echo "Failed to update replica set configuration."
		exit 1
	fi
	echo "Replica set configuration updated."
fi
