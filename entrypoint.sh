#!/bin/sh
set -e

UPLOADS_DIR="${UPLOADS_DIR:-/uploads}"
mkdir -p "$UPLOADS_DIR"
chown -R appuser:appuser "$UPLOADS_DIR"

exec su-exec appuser:appuser "$@"
