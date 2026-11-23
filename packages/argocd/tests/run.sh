#!/usr/bin/env bash
# Render and validate the first-party Argo CD chart package.
set -euo pipefail
PACKAGE_DIR=$(cd "$(dirname "$0")/.." && pwd)
ROOT=$(cd "$PACKAGE_DIR/../.." && pwd)
CHART="$ROOT/charts/argocd"
RENDERED=$(mktemp)
trap 'rm -f "$RENDERED"' EXIT

# dependency update also supports clean runners without a preconfigured Helm repo;
# Chart.yaml and Chart.lock keep the selected upstream version and digest immutable.
helm dependency update "$CHART" >/dev/null
helm template argocd "$CHART" --namespace argocd \
  --api-versions monitoring.coreos.com/v1 \
  --values "$PACKAGE_DIR/values-defaults.yaml" > "$RENDERED"
grep -q 'argocd-server' "$RENDERED"
grep -q 'argocd-repo-server' "$RENDERED"
grep -q 'kind: ServiceMonitor' "$RENDERED"
grep -q 'kind: StatefulSet' "$RENDERED"
! grep -q 'inari-operator' "$RENDERED"
"${HIROBA_PYTHON:-python3}" -c 'import json,sys,jsonschema; schema=json.load(open(sys.argv[1])); values=__import__("yaml").safe_load(open(sys.argv[2])); jsonschema.validate(values,schema)' "$PACKAGE_DIR/values.schema.json" "$PACKAGE_DIR/values-defaults.yaml"
echo "argocd (platform-app): OK"
