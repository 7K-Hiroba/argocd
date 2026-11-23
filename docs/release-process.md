# Release process

Release-please opens one PR with independent components. Do not merge a Release PR or dispatch protected publication without owner/admin approval and exact-head checks. Protected rehearsal, signed tags/Releases, and stable GHCR publication are separate gates; no publication is performed by repository setup or this draft migration.

| Artifact | Repository/ref shape | Identity |
| --- | --- | --- |
| Backend image | `ghcr.io/7k-hiroba/argocd/extension:<version>` | `extension-v<version>` |
| OCI UI remote | `ghcr.io/7k-hiroba/argocd/extension-ui:<version>` | `ui-v<version>` |
| OCI chart | `ghcr.io/7k-hiroba/argocd/charts/argocd:<version>` | `chart-argocd-v<version>` |
| App package | `ghcr.io/7k-hiroba/argocd/catalog/argocd:<version>` | `package-argocd-v<version>` |
| Observability pack | `ghcr.io/7k-hiroba/argocd/catalog/observability-pack:<version>` | `package-observability-pack-v<version>` |
| Scaffold template | `ghcr.io/7k-hiroba/argocd/templates/argocd-application:<version>` | `template-argocd-application-v<version>` |
| Smoke suite | `ghcr.io/7k-hiroba/argocd/smoke/argocd:<chart-version>` | Artifact-bound to chart target |
| Catalog/template indexes | `.../catalog/index:latest`, `.../templates/index:latest` | Discovery-only signed indexes |

Every deployable artifact requires immutable build/push, digest capture, cosign keyless signature, SLSA v1 provenance, release-linked signed test evidence, independent verification, and only then channel/index/finalization. Index `latest` is not a deployment pin. No GitHub-Release UI asset, classic Helm index, or mutable latest/edge deployment is permitted.

The current workflow callers use candidate PR #39 SHA `522c3530c94336249928f5a65a5a9f98c46bf540` and are explicitly marked unreleased. Before protected publication, adopt the exact signed workflow release containing that implementation and rerun affected checks. Do not treat the candidate as released or use an unreleased PR head for final publication.

Chart repair uses the standard `chart-repair.yml` workflow against an existing exact annotated tag. It cannot move channels or rewrite history.
