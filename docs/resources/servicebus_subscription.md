---
subcategory: "Messaging"
page_title: "AzureSim: azuresim_servicebus_subscription"
description: |-
  Manages a simulated Service Bus Subscription.
---

# azuresim_servicebus_subscription

Manages a simulated Service Bus Subscription on a parent Topic.

This resource mimics the behavior of the [`azurerm_servicebus_subscription`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/servicebus_subscription) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_servicebus_topic" "example" {
  name         = "events"
  namespace_id = azuresim_servicebus_namespace.example.id
}

resource "azuresim_servicebus_subscription" "example" {
  name                                  = "audit"
  topic_id                              = azuresim_servicebus_topic.example.id
  max_delivery_count                    = 10
  lock_duration                         = "PT30S"
  default_message_ttl                   = "P14D"
  dead_lettering_on_message_expiration  = true
  enable_batched_operations             = true
  status                                = "Active"
}
```

## Argument Reference

* `name` - (Required, ForceNew) Subscription name.
* `topic_id` - (Required, ForceNew) Parent topic ID.
* `max_delivery_count` - (Required) Max delivery count before dead-letter.
* `lock_duration` - (Optional) Lock duration (ISO 8601).
* `default_message_ttl` - (Optional) Default TTL (ISO 8601).
* `dead_lettering_on_filter_evaluation_error` - (Optional) Dead-letter when a filter errors.
* `dead_lettering_on_message_expiration` - (Optional) Dead-letter on TTL expiry.
* `enable_batched_operations` - (Optional) Enable batched operations.
* `requires_session` - (Optional) Require sessions.
* `forward_to` - (Optional) Auto-forward target queue or topic.
* `forward_dead_lettered_messages_to` - (Optional) Auto-forward dead-letter target.
* `status` - (Optional) `Active` or `Disabled`.

## Attributes Reference

* `id` - The Subscription ID, formed as `<topic_id>/subscriptions/<name>`:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ServiceBus/namespaces/{namespace}/topics/{topic}/subscriptions/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
