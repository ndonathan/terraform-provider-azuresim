---
subcategory: "Key Vault"
page_title: "AzureSim: azuresim_key_vault_secret"
description: |-
  Manages a simulated Azure Key Vault Secret.
---

# azuresim_key_vault_secret

Manages a simulated Azure Key Vault Secret.

This resource mimics the behavior of the [`azurerm_key_vault_secret`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_secret) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_key_vault" "example" {
  name                = "kv-example-001"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_key_vault_secret" "example" {
  name         = "db-password"
  key_vault_id = azuresim_key_vault.example.id
  value        = "P@ssw0rd1234!"
  content_type = "text/plain"

  expiration_date = "2027-12-31T23:59:59Z"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Secret name.
* `key_vault_id` - (Required, ForceNew) Parent Key Vault ID.
* `value` - (Required) Secret value (sensitive).
* `content_type` - (Optional) Free-form content-type tag (e.g. `text/plain`).
* `not_before_date` - (Optional) RFC 3339 timestamp (e.g. `2026-01-01T00:00:00Z`).
* `expiration_date` - (Optional) RFC 3339 timestamp.
* `tags` - (Optional) Tags.

~> **Note:** The `value` is stored in the Terraform state as sensitive data.

## Attributes Reference

* `id` - The versioned data-plane URI:

  ```
  https://<vault-name>.vault.azure.net/secrets/<name>/<version>
  ```

* `version` - Simulated secret version (32-char hex, deterministic from vault, name, and value).
* `versionless_id` - Data-plane URI without version: `https://<vault-name>.vault.azure.net/secrets/<name>`.
* `resource_id` - Resource Manager ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.KeyVault/vaults/{vault_name}/secrets/{name}/{version}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
