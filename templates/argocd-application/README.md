# argocd-application template

Generates a native ArgoCD `Application` and `Kustomization` manifest for review. It requires an explicit Git revision and confined source path; it does not deploy, create a companion, provision observability, enable automatic production sync, or create an AppProject. Use an existing reviewed AppProject and normal GitOps review/sync flow.
