---
subcategory: "Key Vault"
page_title: "AzureSim: azuresim_key_vault_key"
description: |-
  Manages a simulated Azure Key Vault Key.
---

# azuresim_key_vault_key

Manages a simulated Azure Key Vault Key.

This resource mimics the behavior of the [`azurerm_key_vault_key`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_key_vault" "example" {
  name                = "kv-example-001"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_key_vault_key" "example" {
  name         = "data-encryption-key"
  key_vault_id = azuresim_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  expiration_date = "2027-12-31T23:59:59Z"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Key name.
* `key_vault_id` - (Required, ForceNew) Parent Key Vault ID.
* `key_type` - (Required, ForceNew) `EC`, `EC-HSM`, `RSA`, or `RSA-HSM`.
* `key_size` - (Optional) RSA key size (e.g. `2048`, `3072`, `4096`).
* `curve` - (Optional) EC curve (e.g. `P-256`, `P-384`, `P-521`, `P-256K`).
* `key_opts` - (Required) Permitted operations: any of `encrypt`, `decrypt`, `sign`, `verify`, `wrapKey`, `unwrapKey`.
* `not_before_date` - (Optional) RFC 3339 timestamp (e.g. `2026-01-01T00:00:00Z`).
* `expiration_date` - (Optional) RFC 3339 timestamp.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The versioned data-plane URI:

  ```
  https://<vault-name>.vault.azure.net/keys/<name>/<version>
  ```

* `version` - Simulated key version (32-char hex, deterministic).
* `versionless_id` - Data-plane URI without version: `https://<vault-name>.vault.azure.net/keys/<name>`.
* `resource_id` - Resource Manager ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.KeyVault/vaults/{vault_name}/keys/{name}/{version}
  ```

* `public_key_pem` - Placeholder PEM-encoded public key. Not a real key — do not use for any cryptographic operation.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
