#!/usr/bin/env python3
"""Generate the PARITY.md feature-parity report against AzureRM.

Data sources:

* AzureSIM resource list — parsed from ``internal/provider/provider.go``
  by extracting the ``New*Resource`` entries in the ``Resources()``
  slice, then mapping them to ``azuresim_*`` Terraform type names via
  the matching ``resp.TypeName = req.ProviderTypeName + "_..."`` lines
  across ``internal/provider/*.go``.

* AzureRM resource list — read from the committed JSON snapshot at
  ``scripts/azurerm-resources.snapshot.json``. The snapshot was captured
  from the public Terraform Registry API
  (https://registry.terraform.io/v1/providers/hashicorp/azurerm) and is
  intentionally vendored so the generator is fully deterministic and
  safe to run in CI (no network calls).

  **Refresh the snapshot when AzureRM ships a new minor/major release**
  by running ``scripts/refresh-azurerm-snapshot.sh`` and committing the
  result alongside the regenerated ``PARITY.md``.

Modes:
    (default)   write the rendered PARITY.md body to stdout
    --write     replace the block between AUTO-GENERATED markers in
                PARITY.md (creating the file with a header if absent)
    --check     exit 1 if PARITY.md is out of sync with what the
                generator would produce (for CI)

The script is intentionally dependency-free (stdlib only) and
deterministic: identical inputs always produce identical output.

Attribute-level gaps (which AzureRM attributes AzureSIM omits on shared
resources) are **out of scope** for this report — the Registry API does
not expose per-attribute schemas, and scraping the AzureRM source is a
follow-up task tracked in TODO.md.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
PROVIDER_GO = REPO_ROOT / "internal" / "provider" / "provider.go"
PROVIDER_DIR = REPO_ROOT / "internal" / "provider"
AZURERM_SNAPSHOT = REPO_ROOT / "scripts" / "azurerm-resources.snapshot.json"
PARITY_MD = REPO_ROOT / "PARITY.md"

START_MARKER = "<!-- AUTO-GENERATED PARITY START -->"
END_MARKER = "<!-- AUTO-GENERATED PARITY END -->"

# Curated "recommended next-up" picks. Hand-maintained — the Registry
# API doesn't expose per-resource popularity, so the agent's judgment
# fills that gap. Update whenever the AzureSIM resource set changes
# meaningfully or scenario priorities shift.
#
# Format: list of (azurerm_title, one_line_rationale).
RECOMMENDED_NEXT_UP: list[tuple[str, str]] = [
    (
        "linux_virtual_machine",
        "AzureSIM has windows_virtual_machine and the legacy generic "
        "virtual_machine but no Linux split — the most common Azure VM "
        "type in real-world modules.",
    ),
    (
        "lb_backend_address_pool",
        "AzureSIM has azuresim_lb but no children; backend pool is "
        "required for any non-trivial load balancer config.",
    ),
    (
        "lb_rule",
        "Pairs with lb_backend_address_pool and lb_probe to make the LB "
        "resource actually useful end-to-end.",
    ),
    (
        "dns_zone",
        "AzureSIM has private_dns_zone but no public DNS — common for "
        "TLS / cert-manager / ingress demos.",
    ),
    (
        "key_vault_certificate",
        "AzureSIM has key_vault_secret and key_vault_key but not "
        "certificates — required for TLS-binding scenarios.",
    ),
]


def load_azuresim_resources() -> list[str]:
    """Return the sorted list of azuresim_* Terraform type names.

    Cross-checks two sources for safety:
      1. ``New*Resource`` constructor names in provider.go's Resources() slice.
      2. ``resp.TypeName = req.ProviderTypeName + "_..."`` lines across
         every Go file in internal/provider/.

    Raises if the counts disagree — that means a constructor was added
    to the slice but its TypeName wasn't wired (or vice versa).
    """
    provider_text = PROVIDER_GO.read_text(encoding="utf-8")
    constructor_matches = re.findall(
        r"^\s+(New[A-Za-z0-9]+Resource),?\s*$",
        provider_text,
        re.MULTILINE,
    )
    if not constructor_matches:
        raise SystemExit(
            "provider.go: failed to parse any New*Resource entries from the "
            "Resources() slice — has the file structure changed?"
        )

    type_names: list[str] = []
    for go_file in sorted(PROVIDER_DIR.glob("*.go")):
        if go_file.name.endswith("_test.go"):
            continue
        for match in re.finditer(
            r'resp\.TypeName\s*=\s*req\.ProviderTypeName\s*\+\s*"_([a-z0-9_]+)"',
            go_file.read_text(encoding="utf-8"),
        ):
            type_names.append(match.group(1))

    if len(constructor_matches) != len(type_names):
        raise SystemExit(
            f"resource-count mismatch: provider.go declares "
            f"{len(constructor_matches)} New*Resource entries but "
            f"{len(type_names)} TypeName registrations were found in "
            f"internal/provider/*.go. Reconcile the two before regenerating."
        )

    return sorted(set(type_names))


def load_azurerm_resources() -> tuple[list[dict], dict]:
    """Return (resources, snapshot_metadata) from the vendored snapshot."""
    if not AZURERM_SNAPSHOT.is_file():
        raise SystemExit(
            f"missing AzureRM snapshot: {AZURERM_SNAPSHOT}\n"
            "Run scripts/refresh-azurerm-snapshot.sh to create it."
        )
    snap = json.loads(AZURERM_SNAPSHOT.read_text(encoding="utf-8"))
    resources = snap["resources"]
    meta = {
        "source": snap["source"],
        "azurerm_version": snap["azurerm_version"],
        "azurerm_published_at": snap["azurerm_published_at"],
        "captured_at": snap["captured_at"],
    }
    return resources, meta


def render_report(
    azuresim: list[str],
    azurerm: list[dict],
    meta: dict,
) -> str:
    azurerm_titles = {r["title"]: r["subcategory"] for r in azurerm}
    azuresim_set = set(azuresim)

    matched = sorted(azuresim_set & set(azurerm_titles))
    azuresim_only = sorted(azuresim_set - set(azurerm_titles))
    azurerm_only = sorted(set(azurerm_titles) - azuresim_set)

    coverage_pct = (len(matched) / len(azurerm_titles)) * 100 if azurerm_titles else 0.0

    # Group missing AzureRM resources by subcategory for the gap table.
    missing_by_subcat: dict[str, list[str]] = {}
    for title in azurerm_only:
        missing_by_subcat.setdefault(azurerm_titles[title], []).append(title)
    subcat_rows = sorted(
        missing_by_subcat.items(),
        key=lambda kv: (-len(kv[1]), kv[0]),
    )

    lines: list[str] = []
    lines.append(START_MARKER)
    lines.append("")
    lines.append(
        "> This block is auto-generated by `scripts/generate-parity-report.py`. "
        "Do not edit by hand. Run the script after a resource is added or "
        "after refreshing `scripts/azurerm-resources.snapshot.json`; CI "
        "verifies the file stays in sync."
    )
    lines.append("")
    lines.append("## Coverage summary")
    lines.append("")
    lines.append(
        f"**{len(matched)} / {len(azurerm_titles)} AzureRM resources implemented "
        f"({coverage_pct:.1f}%).**"
    )
    lines.append("")
    lines.append(
        f"- AzureSIM-registered resource types: **{len(azuresim_set)}**"
    )
    lines.append(
        f"- AzureRM resource types in snapshot: **{len(azurerm_titles)}**"
    )
    lines.append(
        f"- Name-matched (1:1 `azuresim_X` ↔ `azurerm_X`): **{len(matched)}**"
    )
    if azuresim_only:
        lines.append(
            f"- AzureSIM-only (no AzureRM equivalent by name): "
            f"**{len(azuresim_only)}** — {', '.join('`azuresim_' + n + '`' for n in azuresim_only)}"
        )
    else:
        lines.append(
            "- AzureSIM-only (no AzureRM equivalent by name): **0**"
        )
    lines.append("")
    lines.append(
        f"_AzureRM snapshot: v{meta['azurerm_version']} "
        f"(published {meta['azurerm_published_at']}, "
        f"captured {meta['captured_at']}). "
        f"Source: {meta['source']}._"
    )
    lines.append("")

    # Recommended next-up — curated, agent-judgment-driven.
    lines.append("## Recommended next-up resources")
    lines.append("")
    lines.append(
        "Hand-curated picks for the next round of additions. Updated "
        "alongside `BACKLOG.md` whenever priorities shift; not generated "
        "from the snapshot (the Registry API doesn't expose per-resource "
        "popularity)."
    )
    lines.append("")
    lines.append("| # | AzureRM resource | Rationale |")
    lines.append("|---|---|---|")
    for i, (title, rationale) in enumerate(RECOMMENDED_NEXT_UP, start=1):
        marker = (
            "" if title in azurerm_titles else " _(not in current snapshot — verify name)_"
        )
        lines.append(f"| {i} | `azurerm_{title}`{marker} | {rationale} |")
    lines.append("")

    # Resource-level gaps — grouped by subcategory.
    lines.append("## Resource-level gaps (by AzureRM subcategory)")
    lines.append("")
    lines.append(
        f"AzureRM defines **{len(azurerm_only)}** resources that AzureSIM "
        f"does not yet implement, spread across **{len(missing_by_subcat)}** "
        f"subcategories. Grouped below by subcategory, largest first. "
        f"`azuresim_*` names are inferred — verify naming before adding."
    )
    lines.append("")
    lines.append(
        "<details><summary>Click to expand the full gap list</summary>"
    )
    lines.append("")
    for subcat, titles in subcat_rows:
        lines.append(f"### {subcat} ({len(titles)})")
        lines.append("")
        for title in sorted(titles):
            lines.append(f"- `azurerm_{title}`")
        lines.append("")
    lines.append("</details>")
    lines.append("")

    # Implemented resources — sanity reference.
    lines.append("## Implemented resources")
    lines.append("")
    lines.append(
        f"All {len(matched)} AzureSIM resources map 1:1 by name to an "
        f"AzureRM equivalent. Grouped by AzureRM subcategory for cross-"
        f"reference:"
    )
    lines.append("")
    impl_by_subcat: dict[str, list[str]] = {}
    for title in matched:
        impl_by_subcat.setdefault(azurerm_titles[title], []).append(title)
    for subcat in sorted(impl_by_subcat):
        names = sorted(impl_by_subcat[subcat])
        rendered = ", ".join(f"`azuresim_{n}`" for n in names)
        lines.append(f"- **{subcat}** ({len(names)}): {rendered}")
    lines.append("")

    # Attribute-level gaps — deferred.
    lines.append("## Attribute-level gaps")
    lines.append("")
    lines.append(
        "**Deferred to a follow-up.** The Terraform Registry API exposes "
        "the resource _list_ but not per-resource attribute schemas. "
        "Surfacing AzureRM attributes that AzureSIM omits requires either "
        "(a) parsing each resource's source schema from "
        "`hashicorp/terraform-provider-azurerm/internal/services/*/`, or "
        "(b) introspecting the Registry's per-resource doc-page HTML. "
        "Both are larger pieces of work and are tracked in TODO.md."
    )
    lines.append("")
    lines.append(END_MARKER)
    return "\n".join(lines)


def replace_block(text: str, new_block: str) -> str:
    pattern = re.compile(
        re.escape(START_MARKER) + r".*?" + re.escape(END_MARKER),
        re.DOTALL,
    )
    if not pattern.search(text):
        raise SystemExit(
            "PARITY.md is missing the AUTO-GENERATED markers; "
            "ensure the file has the START/END markers before running --write."
        )
    return pattern.sub(new_block, text)


PARITY_HEADER = """# AzureSIM ↔ AzureRM feature parity

This document tracks where `terraform-provider-azuresim` stands relative
to the upstream `hashicorp/terraform-provider-azurerm`. The block below
is auto-generated — see `scripts/generate-parity-report.py` for how it
is produced and `scripts/azurerm-resources.snapshot.json` for the
vendored ground-truth list.

> **Why a snapshot, not a live API call?** So the parity check is
> deterministic in CI. Refresh the snapshot when AzureRM ships a new
> release that may have added resources by running
> `scripts/refresh-azurerm-snapshot.sh`, then re-run the generator.

"""


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    group = parser.add_mutually_exclusive_group()
    group.add_argument(
        "--write", action="store_true", help="rewrite PARITY.md in place"
    )
    group.add_argument(
        "--check",
        action="store_true",
        help="exit 1 if PARITY.md is out of sync with the generator",
    )
    args = parser.parse_args()

    azuresim = load_azuresim_resources()
    azurerm, meta = load_azurerm_resources()
    new_block = render_report(azuresim, azurerm, meta)

    if args.write:
        if PARITY_MD.is_file():
            existing = PARITY_MD.read_text(encoding="utf-8")
            updated = replace_block(existing, new_block)
        else:
            updated = PARITY_HEADER + new_block + "\n"
        if not PARITY_MD.is_file() or updated != PARITY_MD.read_text(
            encoding="utf-8"
        ):
            PARITY_MD.write_text(updated, encoding="utf-8")
            print(f"updated {PARITY_MD}", file=sys.stderr)
        else:
            print(f"{PARITY_MD} already up to date", file=sys.stderr)
        return 0

    if args.check:
        if not PARITY_MD.is_file():
            print(
                f"{PARITY_MD} does not exist. "
                "Run: python3 scripts/generate-parity-report.py --write",
                file=sys.stderr,
            )
            return 1
        text = PARITY_MD.read_text(encoding="utf-8")
        match = re.search(
            re.escape(START_MARKER) + r".*?" + re.escape(END_MARKER),
            text,
            re.DOTALL,
        )
        if not match:
            print("PARITY.md is missing the AUTO-GENERATED markers", file=sys.stderr)
            return 1
        if match.group(0) != new_block:
            print(
                "PARITY.md is out of sync with the generator. "
                "Run: python3 scripts/generate-parity-report.py --write",
                file=sys.stderr,
            )
            return 1
        return 0

    print(new_block)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
