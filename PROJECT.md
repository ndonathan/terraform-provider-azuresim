---
title: terraform-provider-azuresim
type: terraform-provider
status: active
---

# terraform-provider-azuresim

## Description

A simulated Azure Terraform provider that mimics the AzureRM provider interface — letting you `terraform plan/apply` Azure-style infrastructure with no real API calls, credentials, or cost. Resource state lives entirely in Terraform's state file and computed values (IPs, FQDNs, access keys) are deterministic placeholders. Used for learning, CI/module testing, and demos. Supports 50+ resource types as of the latest BACKLOG.md sweep.

## SDLC compliance

### Universal floor
- [x] README.md
- [x] PROJECT.md
- [ ] TODO.md (BACKLOG.md serves as feature backlog; create TODO.md for non-feature work)
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

## Notes

- Self-published via GitHub Pages registry; consumed by `terraform-iac-azuresim`.
- Signing keys in `gpg-terraform-azuresim/`.
