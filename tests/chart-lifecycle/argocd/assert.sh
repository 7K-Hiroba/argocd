#!/usr/bin/env bash
set -euo pipefail

release=$HIROBA_CHART_RELEASE
namespace=$HIROBA_CHART_NAMESPACE
server=$release-server
repo_server=$release-repo-server
controller=$release-application-controller
actual_digest="sha256:$(sha256sum "$HIROBA_CHART_PACKAGE" | awk '{print $1}')"
[[ "$actual_digest" == "$HIROBA_CHART_DIGEST" ]] || { echo "phase package digest differs from its declared evidence identity" >&2; exit 1; }

# Targeted lifecycle wiring only: the shared runner already verified the Helm
# revision/version. Do not reproduce Argo CD's upstream functional E2E here.
kubectl rollout status "deployment/$server" --namespace "$namespace" --timeout=120s
kubectl rollout status "deployment/$repo_server" --namespace "$namespace" --timeout=120s
kubectl rollout status "statefulset/$controller" --namespace "$namespace" --timeout=120s
kubectl get service "$server" --namespace "$namespace" >/dev/null
kubectl get service "$repo_server" --namespace "$namespace" >/dev/null
server_uid=$(kubectl get service "$server" --namespace "$namespace" -o jsonpath='{.metadata.uid}')
[[ -n "$server_uid" ]] || { echo "server Service UID is empty" >&2; exit 1; }

case "$HIROBA_CHART_PHASE" in
  install)
    printf '%s\n' "$HIROBA_CHART_DIGEST" > "$HIROBA_CHART_WORK/previous-package-digest"
    printf '%s\n' "$server_uid" > "$HIROBA_CHART_WORK/previous-server-service-uid"
    ;;
  upgrade)
    [[ "$HIROBA_CHART_DIGEST" != "$(cat "$HIROBA_CHART_WORK/previous-package-digest")" ]] || { echo "upgrade package did not change content" >&2; exit 1; }
    [[ "$server_uid" == "$(cat "$HIROBA_CHART_WORK/previous-server-service-uid")" ]] || { echo "upgrade replaced the stable server Service identity" >&2; exit 1; }
    ;;
  *) echo "unknown lifecycle assertion phase" >&2; exit 2 ;;
esac
