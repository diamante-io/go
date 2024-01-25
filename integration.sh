#! /bin/bash
set -e

cd "$(dirname "${BASH_SOURCE[0]}")"

export aurora_INTEGRATION_TESTS=true
export aurora_INTEGRATION_ENABLE_CAP_35=${aurora_INTEGRATION_ENABLE_CAP_35:-}
export aurora_INTEGRATION_ENABLE_CAPTIVE_CORE=${aurora_INTEGRATION_ENABLE_CAPTIVE_CORE:-}
export CAPTIVE_CORE_BIN=${CAPTIVE_CORE_BIN:-/usr/bin/diamcircle-core}
export TRACY_NO_INVARIANT_CHECK=1 # This fails on my dev vm. - Paul

# launch postgres if it's not already.
if [[ "$(docker inspect integration_postgres -f '{{.State.Running}}')" != "true" ]]; then
  docker rm -f integration_postgres || true;
  docker run -d --name integration_postgres --env POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 circleci/postgres:9.6.5-alpine
fi

exec go test -timeout 25m go/services/aurora/internal/integration/... "$@"
