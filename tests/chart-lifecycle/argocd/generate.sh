#!/usr/bin/env bash
set -euo pipefail

FIXTURE=$(cd "$(dirname "$0")" && pwd)
ROOT=$(cd "$FIXTURE/../../.." && pwd)
WF=${WF:?WF must identify the exact signed workflows v2.1.2 checkout}
WF=$(cd "$WF" && pwd)
HELM=${HELM:?HELM must identify the checksum-verified pinned Helm binary}
HELM=$(realpath "$(command -v "$HELM")")
PYTHON=${PYTHON:-python3}
EXPECTED_VERSION=$(jq -r .tools.helm.version "$WF/pins/stack-testing.json")
EXPECTED_ARCHIVE=$(jq -r .tools.helm.sha256 "$WF/pins/stack-testing.json")
WORKFLOW_SHA=e105f9c1dd3b305b2141be37ccc773683992f67a
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
mkdir -p "$WORK/helm-plugins" "$WORK/repository-cache"
export HELM_PLUGINS="$WORK/helm-plugins"
export HELM_REPOSITORY_CONFIG="$WORK/repositories.yaml"
export HELM_REPOSITORY_CACHE="$WORK/repository-cache"
export PATH="$(dirname "$HELM"):$PATH"
ACTUAL_VERSION=$("$HELM" version --template '{{.Version}}')
ACTUAL_VERSION=${ACTUAL_VERSION#v}
[[ "$ACTUAL_VERSION" == "$EXPECTED_VERSION" ]] || { echo "Helm $EXPECTED_VERSION is required, found $ACTUAL_VERSION" >&2; exit 2; }

cp -R "$FIXTURE/previous-chart" "$WORK/previous"
cp -R "$ROOT/charts/argocd" "$WORK/current"
# Never reuse an ignored working-tree dependency archive. Resolve both phases
# from their reviewed Chart.yaml/Chart.lock through the signed-R preparer.
rm -rf "$WORK/previous/charts" "$WORK/current/charts"
"$PYTHON" "$WF/scripts/prepare-chart.py" --chart "$WORK/previous"
"$PYTHON" "$WF/scripts/prepare-chart.py" --chart "$WORK/current"
find "$WORK/previous" "$WORK/current" -type f -exec touch -d '@0' {} +
mkdir "$WORK/previous-package" "$WORK/current-package"
PREVIOUS_VERSION=$("$HELM" show chart "$WORK/previous" | awk '/^version:/ {print $2}')
CURRENT_VERSION=$("$HELM" show chart "$WORK/current" | awk '/^version:/ {print $2}')
"$HELM" package "$WORK/previous" --destination "$WORK/previous-package" >/dev/null
"$HELM" package "$WORK/current" --destination "$WORK/current-package" >/dev/null
PYTHONPATH="$WF/scripts" "$PYTHON" -m lib.chart_lifecycle canonicalize-package \
  --input "$WORK/previous-package/argocd-$PREVIOUS_VERSION.tgz" --output "$FIXTURE/previous.tgz"
PYTHONPATH="$WF/scripts" "$PYTHON" -m lib.chart_lifecycle canonicalize-package \
  --input "$WORK/current-package/argocd-$CURRENT_VERSION.tgz" --output "$FIXTURE/current.tgz"
"$PYTHON" - "$FIXTURE" "$EXPECTED_VERSION" "$EXPECTED_ARCHIVE" "$WORKFLOW_SHA" "$PREVIOUS_VERSION" "$CURRENT_VERSION" <<'PY'
import hashlib
import json
import sys
from pathlib import Path

fixture = Path(sys.argv[1])
packages = []
for phase, filename, version, source in (
    ("install", "previous.tgz", sys.argv[5], "tests/chart-lifecycle/argocd/previous-chart"),
    ("upgrade", "current.tgz", sys.argv[6], "charts/argocd"),
):
    digest = hashlib.sha256((fixture / filename).read_bytes()).hexdigest()
    packages.append({"phase": phase, "file": filename, "source": source, "chartVersion": version, "sha256": "sha256:" + digest})
document = {
    "schema": "inari.chart-lifecycle-fixture-provenance.v1",
    "generator": {
        "helmVersion": sys.argv[2],
        "helmArchiveSha256": sys.argv[3],
        "workflowSha": sys.argv[4],
        "sourceFileMtimes": "unix-epoch",
        "archiveCanonicalization": "sorted members, normalized tar metadata, zero gzip mtime",
        "command": "WF=<signed-R-v2.1.2-checkout> PYTHON=<workflow-python> HELM=<checksum-verified-pinned-helm> ./generate.sh",
    },
    "packages": packages,
}
(fixture / "provenance.json").write_text(json.dumps(document, indent=2, sort_keys=True) + "\n", encoding="utf-8")
PY
sha256sum "$FIXTURE/previous.tgz" "$FIXTURE/current.tgz"
