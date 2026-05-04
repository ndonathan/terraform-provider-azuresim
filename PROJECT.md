---
title: terraform-provider-azuresim
type: terraform-provider
status: active
---

# terraform-provider-azuresim

## Description

A simulated Azure Terraform provider that mimics the AzureRM provider interface — letting you `terraform plan/apply` Azure-style infrastructure with no real API calls, credentials, or cost. Resource state lives entirely in Terraform's state file; computed values (IPs, FQDNs, access keys) are deterministic placeholders. Used for learning, CI/module testing, and demos. **Publication to the official Terraform Registry** (currently self-hosted via GitHub Pages registry) is the next major milestone, paired with closing the per-resource doc gap and rotating to a clean-UID signing keypair. Supports 50 resource types as of the latest BACKLOG.md sweep, with feature-parity-with-AzureRM as the ongoing direction.

## Goals

- [ ] **Per-resource docs to 100% coverage** — close the `docs/resources/` gap (currently 13 of 50). Each resource needs a markdown file with attributes, computed values, examples. **Registry publication blocker.** target: 2026-08-31
- [ ] **README "Resource Reference" section trimmed or generated** — currently has hand-written reference for 8 of 50 resources (stale). Either auto-generate from `docs/resources/` or replace with a pointer to the registry listing. target: 2026-09-30
- [ ] **Clean-UID signing keypair for Registry publication** — coordinated with `gpg-terraform-azuresim`'s carried known issue. The "(testing)" UID is fine for the self-hosted registry but reads badly on the official Registry's public listing. Generate a new keypair with a real UID before publication. target: 2026-09-30
- [ ] **Publish to the official Terraform Registry** — follow the publication flow (Registry namespace, signed release via the new keypair, manifest validation, examples). target: 2026-10-31
- [ ] **Feature-parity tracking with AzureRM** — formalize the gap analysis: which AzureRM resources exist that AzureSIM doesn't; which AzureRM attributes exist on shared resources that AzureSIM omits. Could be generated or maintained. Drives BACKLOG.md prioritization going forward. target: 2026-11-30

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
