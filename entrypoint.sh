#!/bin/bash
set -e

sleep 15

touch /tmp/healthy

exec "$@"
