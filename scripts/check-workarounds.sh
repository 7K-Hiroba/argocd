#!/usr/bin/env bash
set -euo pipefail

fail() { echo "forbidden workaround: $*" >&2; exit 1; }
for path in extension/internal/gateway ui/scripts/build-sdk.mjs ui/webpack.config.mjs; do
  [[ ! -e $path ]] || fail "$path still exists"
done
! grep -R -E 'INARI_LEGACY_GATEWAY_TOKEN|INARI_EXTENSION_GATEWAY_TOKEN|INARI_COMPAT_API_VERSION' extension/cmd extension/internal extension/e2e extension.yaml --include='*.go' --include='*.yaml' || fail 'legacy auth or API compatibility knob'
! grep -R -E 'grpc.NewClient|grpc.CallContentSubtype|protojson|ServiceDesc\(\)' extension/cmd extension/internal extension/e2e --include='*.go' || fail 'raw gateway transport'
! grep -R -E '/inari.extensions.v1.AgentGateway|x-inari-extension-token' extension/cmd extension/internal extension/e2e --include='*.go' || fail 'private gateway procedure or shared token metadata'
! grep -R '/api/extensions/' ui/src --include='*.ts' --include='*.tsx' || fail 'manual UI extension invocation'
! grep -R -E 'build-sdk|postinstall|webpack --mode' ui/package.json || fail 'SDK postinstall or custom build'
! grep -R -E 'httpStatusFor|jsonEscape|func reauthError' extension/cmd extension/internal --include='*.go' || fail 'handwritten re-auth mapper'
! grep -R -E 'releases/download|gh release upload' .github ui extension.yaml || fail 'GitHub Release UI remote or upload'
! grep -R -E 'uses: .*@(main|v1|latest)([[:space:]]|$)' .github/workflows || fail 'mutable workflow reference'

echo 'no workaround or mutable release patterns found'
