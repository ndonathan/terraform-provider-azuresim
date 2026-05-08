---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_eventhub"
description: |-
  Manages a simulated Azure Event Hub.
---

# azuresim_eventhub

Manages a simulated Azure Event Hub within a parent Event Hubs Namespace.

This resource mimics the behavior of the [`azurerm_eventhub`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/eventhub) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_eventhub_namespace" "example" {
  name                = "ehns-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"
}

resource "azuresim_eventhub" "example" {
  name              = "eh-events"
  namespace_id      = azuresim_eventhub_namespace.example.id
  partition_count   = 4
  message_retention = 7
  status            = "Active"
}
```

## Example Usage - With Capture

```terraform
resource "azuresim_eventhub" "captured" {
  name              = "eh-captured"
  namespace_id      = azuresim_eventhub_namespace.example.id
  partition_count   = 2
  message_retention = 1

  capture_description {
    enabled             = true
    encoding            = "Avro"
    interval_in_seconds = 300
    size_limit_in_bytes = 314572800

    destination {
      name                = "EventHubArchive.AzureBlockBlob"
      archive_name_format = "{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}"
      blob_container_name = azuresim_storage_container.example.name
      storage_account_id  = azuresim_storage_account.example.id
    }
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Event Hub name.
* `namespace_id` - (Required, ForceNew) Parent Event Hubs Namespace ID.
* `partition_count` - (Required) Number of partitions.
* `message_retention` - (Required) Retention in days.
* `status` - (Optional) `Active`, `Disabled`, or `SendDisabled`.
* `capture_description` - (Optional) Zero or one `capture_description` block as defined below.

---

A `capture_description` block supports:

* `enabled` - (Required) Enable capture.
* `encoding` - (Required) `Avro` or `AvroDeflate`.
* `interval_in_seconds` - (Optional) Capture interval (60-900).
* `size_limit_in_bytes` - (Optional) Capture size limit.
* `skip_empty_archives` - (Optional) Skip empty archives.
* `destination` - (Required) One `destination` block as defined below.

A `destination` block (under `capture_description`) supports:

* `name` - (Required) Destination type (`EventHubArchive.AzureBlockBlob` or `EventHubArchive.AzureDataLake`).
* `archive_name_format` - (Required) Archive path format with placeholders.
* `blob_container_name` - (Required) Target blob container.
* `storage_account_id` - (Required) Target storage account ID.

## Attributes Reference

* `id` - The Event Hub ID, formed as `<namespace_id>/eventhubs/<name>`. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.EventHub/namespaces/{namespace}/eventhubs/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
