---
title: terraform-provider-azuresim
type: terraform-provider
status: active
---

# terraform-provider-azuresim

## Description

A simulated Azure Terraform provider that mimics the AzureRM provider interface — letting you `terraform plan/apply` Azure-style infrastructure with no real API calls, credentials, or cost. Resource state lives entirely in Terraform's state file; computed values (IPs, FQDNs, access keys) are deterministic placeholders. Used for learning, CI/module testing, and demos. **Publication to the official Terraform Registry** (currently self-hosted via GitHub Pages registry) is the next major milestone, paired with closing the per-resource doc gap and rotating to a clean-UID signing keypair. Supports 61 resource types as of the latest BACKLOG.md sweep, with feature-parity-with-AzureRM as the ongoing direction.

## Goals

- [x] **Per-resource docs to 100% coverage** — close the `docs/resources/` gap. Each resource needs a markdown file with attributes, computed values, examples. **Registry publication blocker.** **Done 2026-05-17** — 61 markdown files in `docs/resources/` match the 61 `New*Resource` entries in `internal/provider/provider.go` (48-doc batch `c391c51` 2026-05-07 brought it from 13 to 61).
- [x] **README "Resource Reference" section trimmed or generated** — **Done 2026-05-18** — chose auto-generation: added `scripts/generate-resource-table.py` (stdlib-only, idempotent), replaced both the stale hand-written "Supported Resources" (8 of 61 with `azurerm_*` equivalents) and "Resource Reference" (8 of 61 with per-attribute tables) sections with a single grouped table emitted between `<!-- AUTO-GENERATED RESOURCES START/END -->` markers — 61 resources across 15 subcategories, each row linking to its `docs/resources/*.md`. CI job `docs-table` runs `--check` on every push/PR. Original target was 2026-09-30 — landed early.
- [ ] **Clean-UID signing keypair for Registry publication** — coordinated with `gpg-terraform-azuresim`'s carried known issue. The "(testing)" UID is fine for the self-hosted registry but reads badly on the official Registry's public listing. Generate a new keypair with a real UID before publication. target: 2026-09-30
- [ ] **Publish to the official Terraform Registry** — follow the publication flow (Registry namespace, signed release via the new keypair, manifest validation, examples). target: 2026-10-31
- [x] **Feature-parity tracking with AzureRM** — **Framework + first-pass landed 2026-05-18.** Added `scripts/generate-parity-report.py` (stdlib-only, deterministic) plus vendored `scripts/azurerm-resources.snapshot.json` (Terraform Registry API capture of `hashicorp/azurerm` v4.73.0, refreshed via `scripts/refresh-azurerm-snapshot.sh`). Emits `PARITY.md` between `<!-- AUTO-GENERATED PARITY START/END -->` markers — coverage summary (61 / 1128 = 5.4%), curated next-up table (BACKLOG.md feeder), resource-level gap list grouped by AzureRM subcategory, implemented-resources cross-reference. CI job `parity-table` runs `--check` on push/PR. **Attribute-level gap analysis is deferred to a follow-up** — Registry API doesn't expose per-attribute schemas, so it needs source-clone parsing or HTML scraping; tracked in TODO.md. Original target 2026-11-30 — landed early.

## SDLC compliance

### Universal floor
- [x] README.md
- [x] PROJECT.md
- [x] TODO.md (BACKLOG.md tracks feature work)
- [x] .gitignore

### Project-specific (terraform-provider)
- [x] main.go entry point
- [x] internal/provider/ resources
- [x] docs/
- [x] examples/
- [x] go.mod / go.sum
- [x] BACKLOG.md tracking proposed resources
- [x] .goreleaser.yml for release tooling
- [x] terraform-registry-manifest.json
- [x] Acceptance / unit tests (`internal/provider/*_test.go`)
- [x] CI workflows (`.github/workflows/test.yml`, `release.yml`)
- [x] GNUmakefile (build/test convenience)

## Notes

- Self-published via GitHub Pages registry; consumed by `terraform-iac-azuresim`.
- Signing keys in `gpg-terraform-azuresim/`.
