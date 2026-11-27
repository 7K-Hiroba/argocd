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
- [x] Final workflow SHA is signed R `c2432e1b186987ffdaa97f99b31b89b073c7a47e` (`v2.1.1`, protected run 45665644269), not an open PR head.

## Required acceptance evidence still open

- [x] Hosted PR backend compatibility lanes pass at migration head `31ae2f8` (run 45303387570).
- [x] Nightly backend/browser edge compatibility passed at baseline main `d2fbac0b6d9b45d46cc3961ea71882a22869ce46` with all five slots (run 45613862426; verified artifact 15021861549). Future Release PR heads require fresh evidence.
- [x] Hosted browser slot coverage passes all five slots (run 45303387570).
- [x] Hosted KRO, chart install, and actual Argo CD runtime smoke lanes pass (runs 45303387750, 45303387571 attempt 2, and 45303387656).
- [ ] Stable signed artifacts have cosign/digest/SLSA/test-evidence verification receipts.
- [ ] Exact artifacts admit through trusted source sync and deploy/render/suite admission, including negative checks.
- [ ] SDK-only clean build is re-run at the final head.
- [ ] Sibling Phase 3 is independently confirmed.
- [ ] Central and local companions render through reviewed baselines and target-side verification produces readiness-gated links.
- [ ] Team-configured stage/prod verification demonstrates exact grants, guardrails, dedupe, status/history, and no rollback.
- [ ] `argocd` community registry entry is activated with actual ownership/artifact evidence.
- [ ] SDK task #11 and launch hygiene are rechecked before any public reference promotion.

Nothing in this checklist asserts signed, admitted, rendered, deployed, stage/prod, or production readiness before the corresponding evidence exists.
