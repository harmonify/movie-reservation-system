# https://github.com/elastic/connectors/blob/bc92e92b88951ee6bd5c3263bd3278cbe5a19f2b/scripts/stack/update-kibana-user-password.sh

#!/bin/bash

if [ $# -eq 0 ]; then
	ELASTICSEARCH_URL="http://localhost:9200"
	ELASTIC_PASSWORD="elastic"
else
	ELASTICSEARCH_URL="$1"
	ELASTIC_PASSWORD="$2"
	shift
fi

if [[ ${CURDIR:-} == "" ]]; then
	export CURDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
fi

echo "Updating Kibana password in Elasticsearch running on $ELASTICSEARCH_URL"
change_data="{ \"password\": \"${ELASTIC_PASSWORD}\" }"
curl -u elastic:$ELASTIC_PASSWORD "$@" -X POST "${ELASTICSEARCH_URL}/_security/user/kibana_system/_password?pretty" -H 'Content-Type: application/json' -d"${change_data}"
