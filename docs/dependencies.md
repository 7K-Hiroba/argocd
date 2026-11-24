# Dependency evidence

This repository binds exact dependency commits. **Open candidate PRs are not merged evidence, and neither candidate nor merged-main inputs are signed-release, admitted, rendered, or deployed evidence.** Revalidate the release train and replace candidate pins with exact merged main SHAs as they land; protected publication still requires the approved signed-release SHA.

| Input | Exact SHA/version | Status at migration |
| --- | --- | --- |
| 7K-Hiroba/workflows | `f092eab6bb07de64a2d6f2e35ac4c65fca794f76` / `v2.0.3` | Signed R; protected run 45300302921 and OIDC/gitsign verification passed; signed tag object `77c80969523422d596d39899d8014da8cfc87889` targets R |
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

Workflow callers use exact signed R `f092eab6bb07de64a2d6f2e35ac4c65fca794f76` with `v2.0.3` comments. Protected signing run 45300302921 verified the annotated tag with OIDC/gitsign; tag object `77c80969523422d596d39899d8014da8cfc87889` targets R and `update_major=false` left aliases unchanged. The previously signed v1.1.0 tag and commit remain immutable. No mutable `main`, `v1`, `v2`, `latest`, or `edge` reference is used for deployment.
