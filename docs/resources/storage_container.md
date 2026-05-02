---
subcategory: "Storage"
page_title: "AzureSim: azuresim_storage_container"
description: |-
  Manages a simulated Azure Storage Container.
---

# azuresim_storage_container

Manages a simulated Azure Storage Container (blob container).

This resource mimics the behavior of the [`azurerm_storage_container`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_container) resource.

## Example Usage

```terraform
resource "azuresim_storage_account" "example" {
  name                     = "stexample0001"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_storage_container" "example" {
  name                  = "data"
  storage_account_name  = azuresim_storage_account.example.name
  container_access_type = "private"

  metadata = {
    owner = "platform"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Container name (lowercase, 3-63 chars).
* `storage_account_name` - (Required, ForceNew) Parent Storage Account name.
* `container_access_type` - (Optional, Computed) `blob`, `container`, or `private`. Defaults to `private`.
* `metadata` - (Optional) Container metadata.

## Attributes Reference

* `id` - The data-plane URL of the container:

  ```
  https://{storage_account_name}.blob.core.windows.net/{name}
  ```

* `has_immutability_policy` - Always `false` in this simulator.
* `has_legal_hold` - Always `false` in this simulator.
