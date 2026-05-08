---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_servicebus_queue"
description: |-
  Manages a simulated Service Bus Queue.
---

# azuresim_servicebus_queue

Manages a simulated Service Bus Queue.

This resource mimics the behavior of the [`azurerm_servicebus_queue`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/servicebus_queue) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_servicebus_namespace" "example" {
  name                = "sbns-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"
}

resource "azuresim_servicebus_queue" "example" {
  name                                 = "orders"
  namespace_id                         = azuresim_servicebus_namespace.example.id
  max_size_in_megabytes                = 1024
  lock_duration                        = "PT30S"
  default_message_ttl                  = "P14D"
  max_delivery_count                   = 10
  dead_lettering_on_message_expiration = true
  requires_duplicate_detection         = true
  duplicate_detection_history_time_window = "PT10M"
  status                               = "Active"
}
```

## Argument Reference

* `name` - (Required, ForceNew) Queue name.
* `namespace_id` - (Required, ForceNew) Parent namespace ID.
* `max_size_in_megabytes` - (Optional) Max queue size in MB.
* `lock_duration` - (Optional) ISO 8601 duration (e.g. `PT30S`).
* `requires_duplicate_detection` - (Optional) Enable duplicate detection.
* `requires_session` - (Optional) Enable sessions.
* `dead_lettering_on_message_expiration` - (Optional) Dead-letter messages on TTL expiry.
* `max_delivery_count` - (Optional) Max delivery count before dead-letter.
* `default_message_ttl` - (Optional) Default message TTL (ISO 8601).
* `duplicate_detection_history_time_window` - (Optional) Duplicate detection window (ISO 8601).
* `enable_partitioning` - (Optional) Enable partitioning.
* `enable_express` - (Optional) Enable express queue.
* `status` - (Optional) `Active`, `Disabled`, `SendDisabled`, or `ReceiveDisabled`.

## Attributes Reference

* `id` - The Queue ID, formed as `<namespace_id>/queues/<name>`. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ServiceBus/namespaces/{namespace}/queues/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
