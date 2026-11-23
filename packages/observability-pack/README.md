# observability-pack

Generated v2 KRO companion for Argo CD. Canonical dashboards/rules live in `observability/`; edit those and run `make generate` rather than editing `rgd.yaml`.

The pack renders folder metadata and dashboards for `all`/`dashboards`, central rules for `all`/`rules-central`, and local rules for `all`/`rules-local`, subject to the corresponding boolean toggles. It contains no Keycloak group/team sync, folder ACL, or direct-write component. Target-side verification is required before any readiness link is shown.
