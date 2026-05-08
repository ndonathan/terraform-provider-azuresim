---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_servicebus_namespace"
description: |-
  Manages a simulated Azure Service Bus Namespace.
---

# azuresim_servicebus_namespace

Manages a simulated Azure Service Bus Namespace.

This resource mimics the behavior of the [`azurerm_servicebus_namespace`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/servicebus_namespace) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_servicebus_namespace" "example" {
  name                = "sbns-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"

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
* `capacity` - (Optional) Premium messaging units (1, 2, 4, 8, 16).
* `premium_messaging_partitions` - (Optional) Premium partitions (1, 2, 4).
* `minimum_tls_version` - (Optional) `1.0`, `1.1`, or `1.2`.
* `public_network_access_enabled` - (Optional) Allow public network access.
* `local_auth_enabled` - (Optional) Enable SAS-key authentication.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Service Bus Namespace ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ServiceBus/namespaces/{name}
  ```

* `endpoint` - Simulated SBus endpoint (`sb://<name>.servicebus.windows.net/`).
* `default_primary_key` - Simulated primary key for `RootManageSharedAccessKey` (sensitive).
* `default_secondary_key` - Simulated secondary key (sensitive).
* `default_primary_connection_string` - Simulated primary connection string of the form `Endpoint=sb://<name>.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=<key>` (sensitive).
* `default_secondary_connection_string` - Simulated secondary connection string (sensitive).

~> **Note:** Connection strings and keys are stored in the Terraform state as sensitive data. Values are static placeholders deterministically derived from RG + name and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
