---
subcategory: "Recovery"
page_title: "AzureSim: azuresim_recovery_services_vault"
description: |-
  Manages a simulated Azure Recovery Services Vault.
---

# azuresim_recovery_services_vault

Manages a simulated Azure Recovery Services Vault.

This resource mimics the behavior of the [`azurerm_recovery_services_vault`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/recovery_services_vault) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_recovery_services_vault" "example" {
  name                          = "rsv-example"
  resource_group_name           = azuresim_resource_group.example.name
  location                      = azuresim_resource_group.example.location
  sku                           = "Standard"
  storage_mode_type             = "GeoRedundant"
  soft_delete_enabled           = true
  cross_region_restore_enabled  = true
  immutability                  = "Disabled"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Vault name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Required) `Standard` or `RS0`.
* `storage_mode_type` - (Optional) `GeoRedundant`, `LocallyRedundant`, or `ZoneRedundant`.
* `soft_delete_enabled` - (Optional) Enable soft delete.
* `public_network_access_enabled` - (Optional) Allow public network access.
* `cross_region_restore_enabled` - (Optional) Enable cross-region restore (requires `GeoRedundant`).
* `immutability` - (Optional) `Disabled`, `Unlocked`, or `Locked`.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Vault ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.RecoveryServices/vaults/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
