# Dependency evidence

This repository binds exact dependency commits. **Open candidate PRs are not merged evidence, and neither candidate nor merged-main inputs are signed-release, admitted, rendered, or deployed evidence.** Revalidate the release train and replace candidate pins with exact merged main SHAs as they land; protected publication still requires the approved signed-release SHA.

| Input | Exact SHA/version | Status at migration |
| --- | --- | --- |
| 7K-Hiroba/workflows | `e105f9c1dd3b305b2141be37ccc773683992f67a` / `v2.1.2` | Signed R; protected run 46163314451 attempt 1 and OIDC/gitsign verification passed; signed annotated tag object `fb7f43454063e9d2935969541aa784ec016b9e77` targets R; Release 461682490 pins R with no assets; six same-R rehearsals passed (46060757772, 46078973870, 46079910522, 46082511233, 46085231805, 46143555938); all 28 direct OCI signature/SLSA/predicate checks passed |
| 7K-Hiroba/monitoring | `b1769465757aa0479b18193d47f4a45a60ebd406` | Reviewed candidate PR #2; post-Alloy and readiness inputs, not deployed |
| 7K-Inari/inari-cli | `c79c171dd7b672184087a36697e19e26d1c20559` | Reviewed candidate PR #61; generated the `non-stack` skeleton |
| 7K-Inari/inari-api | `34d8e0db977eb29ac04da8f3eece5224dbcb64a7` / `v0.8.0-34d8e0d` | Merged PR #90 main SHA and exact public tag, consumed by the server candidate for observability profile interoperability |
| 7K-Inari/inari-server | `81ebec882febd60ac150ce9feb8731224b510e57` | Merged PR #230 main SHA for readiness consumption; final main CI 45260731426 green; not a stable release or deployment |
| 7K-Hiroba/community ownership | `d303a92612bfc62b5860c32ca0ab82f971c07965` | Merged PR #13 main SHA; existing hiroba-maintainers ownership input only, with no registry activation |
| 7K-Hiroba/community onboarding | `5651057c0e5489808d5d028f02ffad2b0c585bdb` | Candidate PR #14 |
| inari-plugin-sdk | `v1.2.1-f3f52e4` / `f3f52e4ae28987e04ff25dc8e5e179ba69b411e2` | Exact public Go module tag with typed Agent Gateway and re-auth helpers |
| previous compat SDK | `v1.2.1-e92d21a` / `e92d21a143c21d4fa50b9da1150ff8e6b089c999` | Previous exact public compat tag |
| inari-api Go contracts | `v0.8.0-f224904` | Exact public Go module version required by `inari-plugin-sdk v1.2.1-f3f52e4`; retained until a public SDK tag supports API A90 |
| @7k-inari/ui-plugin-sdk | `0.1.9` | Exact npm registry release with canonical build and host contracts |
| workflow release baseline | `v1.1.0` / `b26b53f5ba58a27d1973fec6dfa6469d0b93ba62` | Previously signed release; does not contain the merged PRs #39/#40/#41 implementation used by current callers |

Workflow callers use exact signed R `e105f9c1dd3b305b2141be37ccc773683992f67a` with `v2.1.2` comments. Protected signing run 46163314451 attempt 1 verified the annotated tag with OIDC/gitsign; tag object `fb7f43454063e9d2935969541aa784ec016b9e77` targets R and `update_major=false` left old tags and aliases unchanged. All six same-R rehearsals and all 28 direct OCI signature/SLSA/predicate checks passed. Previously signed tags and commits, including `v2.1.1` at `c2432e1b186987ffdaa97f99b31b89b073c7a47e`, remain immutable. No mutable `main`, `v1`, `v2`, `latest`, or `edge` reference is used for deployment.
