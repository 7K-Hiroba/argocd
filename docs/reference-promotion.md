# Migration and reference-promotion checklist

## Migration evidence

- [x] Private repository created with existing `hiroba-maintainers` and no new team.
- [x] `non-stack` skeleton generated from exact CLI candidate `c79c171dd7b672184087a36697e19e26d1c20559`.
- [x] Fresh extension import with source repo/path/SHA recorded in `docs/migration/inari-ext-argocd.md`.
- [x] Fresh catalog import with source repo/path/SHA recorded in `docs/migration/inari-catalog.md`.
- [x] No source history, filter-repo, mirror, legacy channel maintenance, or re-pin tooling.
- [x] Backend and UI unit/typecheck/build checks pass; all six workaround families removed.
- [x] Locked OCI chart/package/template/observability/smoke content and generated indexes validated.
- [x] Offline observability and post-Alloy/promtool checks pass.

## Candidate-versus-merged dependency evidence

- [x] Exact candidate SHAs recorded in `docs/dependencies.md`.
- [x] Release train reports exact merged main/signed-release SHAs; candidate pins revalidated/replaced without marking candidates as merged.
- [x] Final workflow SHA is signed R `f092eab6bb07de64a2d6f2e35ac4c65fca794f76` (`v2.0.3`, protected run 45300302921), not an open PR head.

## Required acceptance evidence still open

- [ ] Hosted PR and nightly backend compatibility lanes pass.
- [ ] Hosted browser slot coverage passes all five slots (local harness is pending; see diagnostics).
- [ ] Hosted KRO and actual Argo CD chart runtime lanes pass (local kind environment could not start).
- [ ] Stable signed artifacts have cosign/digest/SLSA/test-evidence verification receipts.
- [ ] Exact artifacts admit through trusted source sync and deploy/render/suite admission, including negative checks.
- [ ] SDK-only clean build is re-run at the final head.
- [ ] Sibling Phase 3 is independently confirmed.
- [ ] Central and local companions render through reviewed baselines and target-side verification produces readiness-gated links.
- [ ] Team-configured stage/prod verification demonstrates exact grants, guardrails, dedupe, status/history, and no rollback.
- [ ] `argocd` community registry entry is activated with actual ownership/artifact evidence.
- [ ] SDK task #11 and launch hygiene are rechecked before any public reference promotion.

Nothing in this checklist asserts signed, admitted, rendered, deployed, stage/prod, or production readiness before the corresponding evidence exists.
