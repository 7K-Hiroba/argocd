# Runtime smoke

This application-evidence-eligible suite deploys the exact `charts/argocd` chart through the explicit non-fixture `helm` adapter, waits for `argocd-server` and `argocd-repo-server`, and performs a read-only Argo CD health and version check through Argo Workflows. It proves an actual Argo CD installation, not a synthetic HTTP fixture.

Only local `ci` executes in this repository. Stage/prod declarations require the separately gated environment verifier and exact approved application/suite digests. Those runs verify the already deployed revision only: they never reinstall, patch desired state, read secrets, create RBAC, or roll back an application. Platform grants bind the exact suite digest and profile to preapproved least-privilege execution settings.
