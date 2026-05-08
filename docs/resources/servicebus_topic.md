---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_servicebus_topic"
description: |-
  Manages a simulated Service Bus Topic.
---

# azuresim_servicebus_topic

Manages a simulated Service Bus Topic.

This resource mimics the behavior of the [`azurerm_servicebus_topic`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/servicebus_topic) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_servicebus_namespace" "example" {
  name                = "sbns-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"
}

resource "azuresim_servicebus_topic" "example" {
  name                                = "events"
  namespace_id                        = azuresim_servicebus_namespace.example.id
  max_size_in_megabytes               = 1024
  default_message_ttl                 = "P14D"
  requires_duplicate_detection        = true
  duplicate_detection_history_time_window = "PT10M"
  support_ordering                    = true
  status                              = "Active"
}
```

## Argument Reference

* `name` - (Required, ForceNew) Topic name.
* `namespace_id` - (Required, ForceNew) Parent namespace ID.
* `max_size_in_megabytes` - (Optional) Max topic size in MB.
* `requires_duplicate_detection` - (Optional) Enable duplicate detection.
* `default_message_ttl` - (Optional) Default message TTL (ISO 8601).
* `duplicate_detection_history_time_window` - (Optional) Duplicate detection window (ISO 8601).
* `enable_partitioning` - (Optional) Enable partitioning.
* `enable_express` - (Optional) Enable express topic.
* `support_ordering` - (Optional) Preserve ordering.
* `status` - (Optional) `Active` or `Disabled`.

## Attributes Reference

* `id` - The Topic ID, formed as `<namespace_id>/topics/<name>`. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ServiceBus/namespaces/{namespace}/topics/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
