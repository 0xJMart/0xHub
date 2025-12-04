#!/bin/sh
set -e
export BACKEND_URL=${BACKEND_URL:-http://0xhub-backend:8080}
envsubst '${BACKEND_URL}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf
exec nginx -g "daemon off;"

