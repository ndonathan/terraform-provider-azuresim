#!/usr/bin/env bash
# Refresh scripts/azurerm-resources.snapshot.json from the public
# Terraform Registry API. Re-run this when the AzureRM provider ships
# a new minor/major version that may have added resources, then commit
# the regenerated JSON alongside the updated PARITY.md.
#
# Usage:  scripts/refresh-azurerm-snapshot.sh
#
# Dependencies: curl, python3 (stdlib only).

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SNAPSHOT="${REPO_ROOT}/scripts/azurerm-resources.snapshot.json"
TMP="$(mktemp -t azurerm-meta.XXXXXX.json)"
trap 'rm -f "${TMP}"' EXIT

echo "Fetching https://registry.terraform.io/v1/providers/hashicorp/azurerm ..." >&2
curl --fail --silent --show-error --max-time 30 \
    "https://registry.terraform.io/v1/providers/hashicorp/azurerm" \
    -o "${TMP}"

python3 - "${TMP}" "${SNAPSHOT}" <<'PY'
import datetime
import json
import sys

src, dst = sys.argv[1], sys.argv[2]
d = json.load(open(src))
resources = [
    {"title": x["title"], "subcategory": x["subcategory"]}
    for x in d["docs"]
    if x["category"] == "resources"
]
resources.sort(key=lambda x: x["title"])
snap = {
    "_comment": (
        "Snapshot of AzureRM provider resources, used by "
        "scripts/generate-parity-report.py. Refresh with: "
        "scripts/refresh-azurerm-snapshot.sh"
    ),
    "source": "https://registry.terraform.io/v1/providers/hashicorp/azurerm",
    "azurerm_version": d["version"],
    "azurerm_published_at": d["published_at"],
    "captured_at": datetime.datetime.now(datetime.timezone.utc).strftime(
        "%Y-%m-%dT%H:%M:%SZ"
    ),
    "resource_count": len(resources),
    "resources": resources,
}
with open(dst, "w", encoding="utf-8") as f:
    json.dump(snap, f, indent=2, ensure_ascii=False)
    f.write("\n")
print(f"wrote {dst} with {len(resources)} resources at AzureRM {snap['azurerm_version']}", file=sys.stderr)
PY

echo "Now run: python3 scripts/generate-parity-report.py --write" >&2
echo "Then review and commit:" >&2
echo "  scripts/azurerm-resources.snapshot.json" >&2
echo "  PARITY.md" >&2
