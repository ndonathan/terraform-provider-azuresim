---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_eventhub_namespace"
description: |-
  Manages a simulated Azure Event Hubs Namespace.
---

# azuresim_eventhub_namespace

Manages a simulated Azure Event Hubs Namespace.

This resource mimics the behavior of the [`azurerm_eventhub_namespace`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/eventhub_namespace) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_eventhub_namespace" "example" {
  name                = "ehns-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"
  capacity            = 1

  auto_inflate_enabled     = true
  maximum_throughput_units = 5
  zone_redundant           = true

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Namespace name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Required) `Basic`, `Standard`, or `Premium`.
* `capacity` - (Optional) Throughput units.
* `auto_inflate_enabled` - (Optional) Enable auto-inflate (`Standard` only).
* `maximum_throughput_units` - (Optional) Auto-inflate ceiling.
* `zone_redundant` - (Optional) Enable zone redundancy.
* `minimum_tls_version` - (Optional) `1.0`, `1.1`, or `1.2`.
* `public_network_access_enabled` - (Optional) Allow public network access.
* `local_auth_enabled` - (Optional) Enable SAS-key authentication.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Event Hubs Namespace ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.EventHub/namespaces/{name}
  ```

* `default_primary_key` - Simulated primary key (sensitive).
* `default_secondary_key` - Simulated secondary key (sensitive).
* `default_primary_connection_string` - Simulated primary connection string of the form `Endpoint=sb://<name>.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=<key>` (sensitive).
* `default_secondary_connection_string` - Simulated secondary connection string (sensitive).

~> **Note:** Connection strings and keys are stored in the Terraform state as sensitive data. Values are static placeholders deterministically derived from RG + name and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
