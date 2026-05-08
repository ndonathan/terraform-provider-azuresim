---
subcategory: "Monitor"
page_title: "AzureSim: azuresim_monitor_metric_alert"
description: |-
  Manages a simulated Azure Monitor metric alert.
---

# azuresim_monitor_metric_alert

Manages a simulated Azure Monitor metric alert.

This resource mimics the behavior of the [`azurerm_monitor_metric_alert`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_metric_alert) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_storage_account" "example" {
  name                     = "stalertexample"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_monitor_action_group" "example" {
  name                = "ag-oncall"
  resource_group_name = azuresim_resource_group.example.name
  short_name          = "oncall"
}

resource "azuresim_monitor_metric_alert" "example" {
  name                = "alert-storage-availability"
  resource_group_name = azuresim_resource_group.example.name
  description         = "Alert when storage account availability drops"
  scopes              = [azuresim_storage_account.example.id]
  enabled             = true
  severity            = 2
  frequency           = "PT1M"
  window_size         = "PT5M"
  auto_mitigate       = true

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "Availability"
    aggregation      = "Average"
    operator         = "LessThan"
    threshold        = 99.9
  }

  action {
    action_group_id = azuresim_monitor_action_group.example.id
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Alert name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `description` - (Optional) Alert description.
* `scopes` - (Required) List of resource IDs to monitor.
* `enabled` - (Optional) Whether the alert is enabled.
* `severity` - (Optional) Severity (0-4).
* `frequency` - (Optional) Evaluation frequency, ISO 8601 (e.g. `PT1M`).
* `window_size` - (Optional) Evaluation window, ISO 8601 (e.g. `PT5M`).
* `auto_mitigate` - (Optional) Auto-mitigate when the condition clears.
* `target_resource_type` - (Optional) Target resource type (for multi-resource alerts).
* `target_resource_location` - (Optional) Target resource region (for multi-resource alerts).
* `criteria` - (Required) One or more `criteria` blocks as defined below.
* `action` - (Optional) Zero or more `action` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `criteria` block supports:

* `metric_namespace` - (Required) Metric namespace.
* `metric_name` - (Required) Metric name.
* `aggregation` - (Required) `Average`, `Count`, `Minimum`, `Maximum`, or `Total`.
* `operator` - (Required) `Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, or `LessThanOrEqual`.
* `threshold` - (Required) Threshold value.
* `skip_metric_validation` - (Optional) Skip Azure validation of the metric name.

---

An `action` block supports:

* `action_group_id` - (Required) Action Group ID.
* `webhook_properties` - (Optional) Custom webhook properties (map of strings).

## Attributes Reference

* `id` - The Metric Alert ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/microsoft.insights/metricAlerts/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
