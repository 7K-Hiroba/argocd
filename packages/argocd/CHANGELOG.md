# Changelog

## [1.4.0](https://github.com/7K-Hiroba/argocd/compare/package-argocd-v1.3.0...package-argocd-v1.4.0) (2026-11-28)


### Features

* add ArgoCD OCI content and verification assets ([fc42a10](https://github.com/7K-Hiroba/argocd/commit/fc42a10e703ddcff13d4db4eade4935c8b6e7d80))
* migrate ArgoCD to the standalone application architecture ([fbb37d3](https://github.com/7K-Hiroba/argocd/commit/fbb37d35d796d6943fbcbdf9c4a9ae4b6fff34af))


### Bug Fixes

* address migration review findings ([56f5b93](https://github.com/7K-Hiroba/argocd/commit/56f5b93c5ea022a029bae21158ef07255bde520f))
* adopt merged workflows and repair chart CI ([525ca46](https://github.com/7K-Hiroba/argocd/commit/525ca46408bc8c6aa310b98b6e12b3520a35d3b5))
* register chart repo for locked tests ([d7a4efe](https://github.com/7K-Hiroba/argocd/commit/d7a4efe020a3970f861603a3cdbbca9449e4291d))

## [1.3.0](https://github.com/7K-Inari/inari-catalog/compare/v1.2.0...v1.3.0) (2026-09-25)


### Features

* M3-W1 Crossplane AWS packages, platform apps, tenant-zone-aws RGD, policy packs ([ef2c72a](https://github.com/7K-Inari/inari-catalog/commit/ef2c72a9817937b735948a010a751acf4bbffa30))
* **platform-apps:** keycloak, cert-manager, external-secrets, argocd, monitoring-stack ([7a198dd](https://github.com/7K-Inari/inari-catalog/commit/7a198dd0cdee98bb5e764e8e92d588f45671c62d))
* publish remaining catalog packages to GHCR ([aa0f143](https://github.com/7K-Inari/inari-catalog/commit/aa0f1435b5a3d73c69cca5fe9c14deff96c68f36))
* publish remaining catalog packages to GHCR ([4d13d4f](https://github.com/7K-Inari/inari-catalog/commit/4d13d4f571caa730329552e5cfc612d83be84c6d))


### Bug Fixes

* **qa:** broken policy CEL, chained-ignore ServiceAccount, unguarded optional lists, helm test assertions ([9d8b6e9](https://github.com/7K-Inari/inari-catalog/commit/9d8b6e901fc7d0533931bbc432e4f1f7a8dc118b))

## [1.2.0](https://github.com/7K-Inari/inari-catalog/compare/v1.1.0...v1.2.0) (2026-09-13)


### Features

* M3-W1 Crossplane AWS packages, platform apps, tenant-zone-aws RGD, policy packs ([ef2c72a](https://github.com/7K-Inari/inari-catalog/commit/ef2c72a9817937b735948a010a751acf4bbffa30))
* **platform-apps:** keycloak, cert-manager, external-secrets, argocd, monitoring-stack ([7a198dd](https://github.com/7K-Inari/inari-catalog/commit/7a198dd0cdee98bb5e764e8e92d588f45671c62d))
* publish remaining catalog packages to GHCR ([aa0f143](https://github.com/7K-Inari/inari-catalog/commit/aa0f1435b5a3d73c69cca5fe9c14deff96c68f36))
* publish remaining catalog packages to GHCR ([4d13d4f](https://github.com/7K-Inari/inari-catalog/commit/4d13d4f571caa730329552e5cfc612d83be84c6d))


### Bug Fixes

* **qa:** broken policy CEL, chained-ignore ServiceAccount, unguarded optional lists, helm test assertions ([9d8b6e9](https://github.com/7K-Inari/inari-catalog/commit/9d8b6e901fc7d0533931bbc432e4f1f7a8dc118b))

## [1.1.0](https://github.com/7K-Inari/inari-catalog/compare/v1.0.0...v1.1.0) (2026-09-13)


### Features

* M3-W1 Crossplane AWS packages, platform apps, tenant-zone-aws RGD, policy packs ([ef2c72a](https://github.com/7K-Inari/inari-catalog/commit/ef2c72a9817937b735948a010a751acf4bbffa30))
* **platform-apps:** keycloak, cert-manager, external-secrets, argocd, monitoring-stack ([7a198dd](https://github.com/7K-Inari/inari-catalog/commit/7a198dd0cdee98bb5e764e8e92d588f45671c62d))
* publish remaining catalog packages to GHCR ([aa0f143](https://github.com/7K-Inari/inari-catalog/commit/aa0f1435b5a3d73c69cca5fe9c14deff96c68f36))
* publish remaining catalog packages to GHCR ([4d13d4f](https://github.com/7K-Inari/inari-catalog/commit/4d13d4f571caa730329552e5cfc612d83be84c6d))


### Bug Fixes

* **qa:** broken policy CEL, chained-ignore ServiceAccount, unguarded optional lists, helm test assertions ([9d8b6e9](https://github.com/7K-Inari/inari-catalog/commit/9d8b6e901fc7d0533931bbc432e4f1f7a8dc118b))

## 1.0.0 (unreleased)

### Features

* argocd platform app: argo-cd chart with HA defaults, redis-ha, metrics and
  ServiceMonitors; dex/notifications disabled (Keycloak SSO).
