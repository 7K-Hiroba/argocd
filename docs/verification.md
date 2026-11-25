# Verification

The `tests/smoke` suite binds `charts/argocd` and uses the non-fixture `helm` adapter. CI installs the exact chart and waits for Argo CD server/repo-server before a read-only health/version request. This is application evidence, not fixture-only evidence.

Stage/prod declarations do not authorize execution. After sibling Phase 3, the generic verifier must import the exact signed suite and bind an exact digest/profile grant to the existing instance revision. Platform policy supplies the approved ServiceAccount, RBAC, egress, concurrency, retention, and expiry. Teams may then configure revision/manual/scheduled triggers, notifications, and promotion gates. Deployed runs never reinstall, patch desired state, read secrets, create RBAC, or roll back.

Expected evidence includes app instance UID/revision, suite digest, policy/grant revision, deterministic run ID/attempt, Workflow UID, outbound agent result, notifications, and promotion decision. Duplicate triggers must not create duplicate runs; unrelated syncs must not trigger; one active run is the default. Failure only notifies, pauses, or blocks promotion as configured. There is no automatic rollback.

The current server/API contract has no ArgoCD rollback history ID. The UI therefore proves safe no-recorded-revision validation in browser coverage. No revision is invented from `observedRevision`; any deployed rollback-history contract extension is a separate parent-coordinated defect.
