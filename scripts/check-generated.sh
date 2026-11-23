#!/usr/bin/env bash
set -euo pipefail
WF=${WF:?WF must point to the exact workflow checkout}
PYTHON=${PYTHON:-python3}
"$PYTHON" "$WF/scripts/render-observability.py" --root . --check
"$PYTHON" "$WF/scripts/build-indexes.py" --root . --check
git diff --exit-code
git diff --check
