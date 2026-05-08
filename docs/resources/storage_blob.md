---
subcategory: "Storage"
page_title: "AzureSim: azuresim_storage_blob"
description: |-
  Manages a simulated Azure Storage Blob.
---

# azuresim_storage_blob

Manages a simulated Azure Storage Blob within a parent Storage Container.

This resource mimics the behavior of the [`azurerm_storage_blob`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_blob) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_storage_account" "example" {
  name                     = "stblobexample"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_storage_container" "example" {
  name                  = "data"
  storage_account_name  = azuresim_storage_account.example.name
  container_access_type = "private"
}

resource "azuresim_storage_blob" "example" {
  name                   = "config.json"
  storage_account_name   = azuresim_storage_account.example.name
  storage_container_name = azuresim_storage_container.example.name
  type                   = "BlockBlob"
  content_type           = "application/json"
  source_content         = jsonencode({ environment = "production" })
  access_tier            = "Hot"

  metadata = {
    owner = "platform"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Blob name.
* `storage_account_name` - (Required, ForceNew) Parent Storage Account name.
* `storage_container_name` - (Required, ForceNew) Parent container name.
* `type` - (Required) `BlockBlob`, `PageBlob`, or `AppendBlob`.
* `size` - (Optional) Blob size in bytes (PageBlob only).
* `content_type` - (Optional) Content-Type header.
* `content_md5` - (Optional) MD5 hash.
* `source` - (Optional) Local file path.
* `source_content` - (Optional) Inline blob content.
* `source_uri` - (Optional) URI to copy from.
* `access_tier` - (Optional) `Hot`, `Cool`, `Cold`, or `Archive`.
* `cache_control` - (Optional) Cache-Control header.
* `metadata` - (Optional) Map of blob metadata.

## Attributes Reference

* `id` - The blob URL (used as the resource ID):

  ```
  https://<storage_account_name>.blob.core.windows.net/<storage_container_name>/<name>
  ```

* `url` - The public blob URL (same as `id`).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
