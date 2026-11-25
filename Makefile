WF ?= .hiroba/workflows
PYTHON ?= python3

.PHONY: test test-backend test-ui validate generate check-generated check-workarounds smoke-ci

test: test-backend test-ui validate check-generated check-workarounds

test-backend:
	cd extension && test -z "$$(gofmt -l .)" && go vet ./... && go test ./... && go build ./... && go mod tidy -diff

test-ui:
	cd ui && npm ci --include=dev --ignore-scripts && npm test && npm run typecheck && npm run lint && npm run build

validate:
	$(PYTHON) $(WF)/scripts/validate-repository.py --root .

generate:
	$(PYTHON) $(WF)/scripts/render-observability.py --root .
	$(PYTHON) $(WF)/scripts/build-indexes.py --root .

check-generated:
	$(PYTHON) $(WF)/scripts/render-observability.py --root . --check
	$(PYTHON) $(WF)/scripts/build-indexes.py --root . --check
	git diff --exit-code
	git diff --check

check-workarounds:
	bash scripts/check-workarounds.sh

smoke-ci:
	HIROBA_PYTHON=$(PYTHON) bash $(WF)/scripts/run-runtime-smoke.sh --root . --mode local --profile ci
