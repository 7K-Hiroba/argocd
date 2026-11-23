# Dependency evidence

This repository binds exact dependency commits. **Open candidate PRs are not merged evidence, and neither candidate nor merged-main inputs are signed-release, admitted, rendered, or deployed evidence.** Revalidate the release train and replace candidate pins with exact merged main SHAs as they land; protected publication still requires the approved signed-release SHA.

| Input | Exact SHA/version | Status at migration |
| --- | --- | --- |
| 7K-Hiroba/workflows | `7e61f10f1f78ebf9470a6f9343dd86d71e4e6447` | Merged PR #39 main SHA, unreleased; includes app-config-v2, Mimir generation, and locked chart dependency preparation |
| 7K-Hiroba/monitoring | `b1769465757aa0479b18193d47f4a45a60ebd406` | Reviewed candidate PR #2; post-Alloy and readiness inputs, not deployed |
| 7K-Inari/inari-cli | `c79c171dd7b672184087a36697e19e26d1c20559` | Reviewed candidate PR #61; generated the `non-stack` skeleton |
| 7K-Inari/inari-api | `734f6d7a10689ec97ff8439f75709e94fc7ee0cf` | Candidate PR #90 for observability profile interoperability |
| 7K-Inari/inari-server | `5feabbd042f8b11356f9016a5e592e0ca008dd6e` | Candidate PR #230 for readiness consumption |
| 7K-Hiroba/community ownership | `3f80974651e0940b1c923405eba38f1fd734d4d6` | Candidate PR #13 |
| 7K-Hiroba/community onboarding | `5651057c0e5489808d5d028f02ffad2b0c585bdb` | Candidate PR #14 |
| inari-plugin-sdk | `v1.2.1-f3f52e4` / `f3f52e4ae28987e04ff25dc8e5e179ba69b411e2` | Exact public Go module tag with typed Agent Gateway and re-auth helpers |
| previous compat SDK | `v1.2.1-e92d21a` / `e92d21a143c21d4fa50b9da1150ff8e6b089c999` | Previous exact public compat tag |
| inari-api Go contracts | `v0.8.0-f224904` | Exact public Go module version |
| @7k-inari/ui-plugin-sdk | `0.1.9` | Exact npm registry release with canonical build and host contracts |
| workflow release baseline | `v1.1.0` / `b26b53f5ba58a27d1973fec6dfa6469d0b93ba62` | Previously signed release; does not contain the merged PR #39 implementation used by current callers |

Workflow callers carry an explicit `v1.1.0+main.7e61f10 (unreleased); merged workflows PR #39` comment so the merged main SHA cannot be mistaken for a signed release. The previously signed v1.1.0 tag and commit remain immutable. Final protected publication still requires a separately approved signed workflow release containing these changes. No mutable `main`, `v1`, `latest`, or `edge` reference is used for deployment.
