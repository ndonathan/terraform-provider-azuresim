---
subcategory: "Storage"
page_title: "AzureSim: azuresim_storage_account"
description: |-
  Manages a simulated Azure Storage Account.
---

# azuresim_storage_account

Manages a simulated Azure Storage Account.

This resource mimics the behavior of the [`azurerm_storage_account`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account) resource. It manages all state within Terraform's state file and does not make any API calls to Azure. Computed attributes such as `primary_blob_endpoint` and `primary_access_key` are generated locally based on the provided configuration.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_storage_account" "example" {
  name                     = "stexampleaccount"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
```

## Example Usage - Premium Block Blob Storage

```terraform
resource "azuresim_storage_account" "premium" {
  name                     = "stpremiumblobs"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "BlockBlobStorage"

  tags = {
    environment = "production"
    tier        = "premium"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Storage Account. In a real Azure environment, this must be globally unique, between 3 and 24 characters in length, and may contain only lowercase letters and numbers. Changing this forces a new Storage Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Storage Account should exist. Changing this forces a new Storage Account to be created.

* `location` - (Required) The Azure Region where the Storage Account should exist. Changing this forces a new Storage Account to be created.

* `account_tier` - (Required) Defines the Tier to use for this Storage Account. Valid options are `Standard` and `Premium`. In a real Azure environment, `Premium` is required for certain storage kinds such as `BlockBlobStorage` and `FileStorage`.

* `account_replication_type` - (Required) Defines the type of replication to use for this Storage Account. Valid options include:
  * `LRS` - Locally-redundant storage.
  * `GRS` - Geo-redundant storage.
  * `RAGRS` - Read-access geo-redundant storage.
  * `ZRS` - Zone-redundant storage.

* `account_kind` - (Optional) Defines the Kind of Storage Account. Valid options are `BlobStorage`, `BlockBlobStorage`, `FileStorage`, `Storage`, and `StorageV2`. Defaults to `StorageV2` in a real Azure environment.

* `tags` - (Optional) A mapping of tags which should be assigned to the Storage Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Storage/storageAccounts/{name}
  ```

* `primary_blob_endpoint` - The endpoint URL for blob storage in the primary location. Generated in the format:

  ```
  https://{name}.blob.core.windows.net/
  ```

* `primary_access_key` - A simulated primary access key for the Storage Account.

~> **Note:** The `primary_access_key` is stored in the Terraform state as sensitive data. The value is a static placeholder and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
