---
subcategory: "Log Analytics"
page_title: "AzureSim: azuresim_log_analytics_workspace"
description: |-
  Manages a simulated Azure Log Analytics Workspace.
---

# azuresim_log_analytics_workspace

Manages a simulated Azure Log Analytics Workspace.

This resource mimics the behavior of the [`azurerm_log_analytics_workspace`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/log_analytics_workspace) resource.

## Example Usage

```terraform
resource "azuresim_log_analytics_workspace" "example" {
  name                = "law-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "PerGB2018"
  retention_in_days   = 90

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the workspace.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Optional, Computed) SKU. One of `Free`, `Standalone`, `PerNode`, `PerGB2018`, `CapacityReservation`. Defaults to `PerGB2018`.
* `retention_in_days` - (Optional, Computed) Retention in days (30-730). Defaults to `30`.
* `daily_quota_gb` - (Optional) Daily ingestion cap in GB. Use `-1` for unlimited.
* `internet_ingestion_enabled` - (Optional) Whether public internet ingestion is allowed.
* `internet_query_enabled` - (Optional) Whether public internet query is allowed.
* `reservation_capacity_in_gb_per_day` - (Optional) Reserved capacity (only with `CapacityReservation` SKU).
* `local_authentication_disabled` - (Optional) Disable non-AAD authentication.
* `cmk_for_query_forced` - (Optional) Force CMK for query.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Workspace resource ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.OperationalInsights/workspaces/{name}
  ```

* `workspace_id` - Simulated Customer ID (UUID), deterministically derived from RG + name.
* `primary_shared_key` - Simulated primary shared key (sensitive).
* `secondary_shared_key` - Simulated secondary shared key (sensitive).
