# ArgoCD backend extension

Go backend for ArgoCD sync, refresh, rollback, and custom resource actions. It runs in host-supervised subprocess mode or HTTP dial mode and uses only public Inari SDK contracts.

## Runtime contract

- `cmd/inari-ext-argocd`: go-plugin subprocess entrypoint.
- `cmd/inari-ext-argocd-http`: ConnectRPC and `/actions/<name>` HTTP dial-mode entrypoint.
- `internal/extplugin`: action registration, strict schemas, auth, health, and typed re-auth.
- `internal/argocd`: payload validation and command construction.
- `internal/agentstub`: tenant-local ArgoCD REST executor used by integration tests.
- `internal/runtimeconfig`: shared public Agent Gateway/OAuth configuration.
- `e2e/`: generated Connect test gateway, round trips, and real kind journey.

Every action requires an Inari-managed AppProject (default `inari`) and a per-user ArgoCD SSO session. Tenant identity comes from `AuthContext`; disconnected agents fail closed. Raw downstream tokens never enter persisted commands, logs, audits, or errors.

## Configuration

- `INARI_AGENT_GATEWAY_URL`: absolute control-plane HTTP(S) URL. Unset means actions and health fail closed while handshake/capabilities remain available.
- `INARI_EXTENSION_CLIENT_ID`, `INARI_EXTENSION_CLIENT_SECRET`, `INARI_EXTENSION_TOKEN_URL`: dedicated extension client-credentials configuration returned by registration or rotation and supplied through deployment secret machinery. Required when the gateway URL is set.
- `INARI_MANAGED_PROJECTS`: comma-separated allowlist (default `inari`).
- `INARI_ACTION_TIMEOUT`: per-action bound (default `60s`).
- `INARI_HTTP_ADDR`: HTTP dial-mode listen address (default `:8080`).

No shared gateway token, raw gateway proto, custom codec, legacy API-version override, or handwritten HTTP re-auth mapping is supported.

## Develop

```sh
go test ./...
go vet ./...
go build ./...
./e2e/kind/run.sh  # optional real kind journey with Docker, kind, and kubectl
```

The E2E harness uses generated `extensionsv1connect` handlers and the typed SDK client, including an ordinary OAuth client-credentials test flow. Public dependency versions are pinned in `go.mod`; no replace/workspace/private contract fork is allowed.
