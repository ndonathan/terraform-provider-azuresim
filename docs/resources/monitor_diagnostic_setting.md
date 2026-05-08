---
subcategory: "Monitor"
page_title: "AzureSim: azuresim_monitor_diagnostic_setting"
description: |-
  Manages a simulated Azure Monitor Diagnostic Setting.
---

# azuresim_monitor_diagnostic_setting

Manages a simulated Azure Monitor Diagnostic Setting on a target resource.

This resource mimics the behavior of the [`azurerm_monitor_diagnostic_setting`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_diagnostic_setting) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_log_analytics_workspace" "example" {
  name                = "law-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "PerGB2018"
}

resource "azuresim_key_vault" "example" {
  name                = "kv-example-001"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_monitor_diagnostic_setting" "example" {
  name                       = "diag-kv"
  target_resource_id         = azuresim_key_vault.example.id
  log_analytics_workspace_id = azuresim_log_analytics_workspace.example.id

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Diagnostic Setting name.
* `target_resource_id` - (Required, ForceNew) ID of the resource the diagnostic setting applies to.
* `log_analytics_workspace_id` - (Optional) Log Analytics Workspace destination ID.
* `storage_account_id` - (Optional) Storage Account destination ID.
* `eventhub_authorization_rule_id` - (Optional) Event Hub authorization rule ID destination.
* `eventhub_name` - (Optional) Event Hub name destination.
* `log_analytics_destination_type` - (Optional) `Dedicated` or `AzureDiagnostics`.
* `enabled_log` - (Optional) Zero or more `enabled_log` blocks as defined below.
* `metric` - (Optional) Zero or more `metric` blocks as defined below.

---

An `enabled_log` block supports:

* `category` - (Optional) Specific log category name.
* `category_group` - (Optional) Log category group (`allLogs`, `audit`, etc.).

---

A `metric` block supports:

* `category` - (Required) Metric category (commonly `AllMetrics`).
* `enabled` - (Optional) Whether this metric category is enabled.

## Attributes Reference

* `id` - The Diagnostic Setting ID, formed as `<target_resource_id>|<name>`.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
