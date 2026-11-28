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

The current workflow callers use exact signed R `e105f9c1dd3b305b2141be37ccc773683992f67a` (`v2.1.2`). Protected signing run 46163314451 attempt 1 verified the release with OIDC/gitsign; annotated tag object `fb7f43454063e9d2935969541aa784ec016b9e77` targets R. Release 461682490 pins R and has no assets. All six same-R rehearsals passed: 46060757772, 46078973870, 46079910522, 46082511233, 46085231805, and 46143555938; all 28 direct OCI signature/SLSA/predicate checks passed. This pin remains fixed even if workflows main advances; no mutable main/tag is substituted. Protected ArgoCD publication remains a separate owner/admin gate and is not performed by this adoption.

Chart repair uses the standard `chart-repair.yml` workflow against an existing exact annotated tag. It cannot move channels or rewrite history.
