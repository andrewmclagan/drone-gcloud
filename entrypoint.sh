#!/bin/sh

set -e

echo $PLUGIN_SERVICE_KEY | base64 -d > service_key.json

CLIENT_EMAIL=$(cat service_key.json | jq -r '.client_email')

PROJECT_ID=$(cat service_key.json | jq -r '.project_id')

gcloud auth activate-service-account \
    $CLIENT_EMAIL \
    --key-file=service_key.json \
    --project=$PROJECT_ID

exec "$@"    