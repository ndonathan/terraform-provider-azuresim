---
subcategory: "Monitor"
page_title: "AzureSim: azuresim_monitor_action_group"
description: |-
  Manages a simulated Azure Monitor Action Group.
---

# azuresim_monitor_action_group

Manages a simulated Azure Monitor Action Group.

This resource mimics the behavior of the [`azurerm_monitor_action_group`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_action_group) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_monitor_action_group" "example" {
  name                = "ag-oncall"
  resource_group_name = azuresim_resource_group.example.name
  short_name          = "oncall"
  enabled             = true

  email_receiver {
    name                    = "oncall-engineer"
    email_address           = "oncall@example.com"
    use_common_alert_schema = true
  }

  sms_receiver {
    name         = "oncall-sms"
    country_code = "1"
    phone_number = "5555550100"
  }

  webhook_receiver {
    name                    = "alerts-bot"
    service_uri             = "https://example.com/alerts"
    use_common_alert_schema = true
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Action Group name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `short_name` - (Required) Short name (max 12 chars).
* `enabled` - (Optional) Whether the action group is enabled.
* `email_receiver` - (Optional) Zero or more `email_receiver` blocks as defined below.
* `sms_receiver` - (Optional) Zero or more `sms_receiver` blocks as defined below.
* `webhook_receiver` - (Optional) Zero or more `webhook_receiver` blocks as defined below.
* `tags` - (Optional) Tags.

---

An `email_receiver` block supports:

* `name` - (Required) Receiver name.
* `email_address` - (Required) Email address.
* `use_common_alert_schema` - (Optional) Use the common alert schema.

---

An `sms_receiver` block supports:

* `name` - (Required) Receiver name.
* `country_code` - (Required) Country code.
* `phone_number` - (Required) Phone number.

---

A `webhook_receiver` block supports:

* `name` - (Required) Receiver name.
* `service_uri` - (Required) Webhook URL.
* `use_common_alert_schema` - (Optional) Use the common alert schema.

## Attributes Reference

* `id` - The Action Group ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/microsoft.insights/actionGroups/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
