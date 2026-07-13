#!/bin/sh
set -e

POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-postgres}"
POSTGRES_DB="${POSTGRES_DB:-markdownkb}"

cat > kubernetes/secret.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
type: Opaque
data:
  POSTGRES_USER: $(echo -n "$POSTGRES_USER" | base64)
  POSTGRES_PASSWORD: $(echo -n "$POSTGRES_PASSWORD" | base64)
  POSTGRES_DB: $(echo -n "$POSTGRES_DB" | base64)
  DATABASE_URL: $(echo -n "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/$POSTGRES_DB?sslmode=disable" | base64)
EOF

echo "Generated kubernetes/secret.yaml"
