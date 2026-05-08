---
subcategory: "Monitor"
page_title: "AzureSim: azuresim_application_insights"
description: |-
  Manages a simulated Azure Application Insights component.
---

# azuresim_application_insights

Manages a simulated Azure Application Insights component.

This resource mimics the behavior of the [`azurerm_application_insights`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/application_insights) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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
  retention_in_days   = 90
}

resource "azuresim_application_insights" "example" {
  name                = "appi-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  application_type    = "web"
  workspace_id        = azuresim_log_analytics_workspace.example.id
  retention_in_days   = 90
  sampling_percentage = 100

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the component.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `application_type` - (Required, ForceNew) `web`, `other`, `java`, `MobileCenter`, `Node.JS`, `phone`, or `store`.
* `workspace_id` - (Optional) Linked Log Analytics Workspace ID for workspace-based instances.
* `retention_in_days` - (Optional) Data retention in days. One of `30`, `60`, `90`, `120`, `180`, `270`, `365`, `550`, `730`.
* `sampling_percentage` - (Optional) Sampling percentage (0-100).
* `daily_data_cap_in_gb` - (Optional) Daily data cap in GB.
* `daily_data_cap_notifications_disabled` - (Optional) Disable email notifications when the daily cap is hit.
* `disable_ip_masking` - (Optional) Disable IP masking on incoming telemetry.
* `local_authentication_disabled` - (Optional) Disable non-AAD authentication.
* `internet_ingestion_enabled` - (Optional) Allow public internet ingestion.
* `internet_query_enabled` - (Optional) Allow public internet query.
* `force_customer_storage_for_profiler` - (Optional) Force customer-managed storage for Profiler.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Application Insights ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Insights/components/{name}
  ```

* `app_id` - Simulated Application ID (UUID), deterministically derived from RG + name.
* `instrumentation_key` - Simulated instrumentation key (UUID, sensitive).
* `connection_string` - Simulated connection string of the form `InstrumentationKey=...;IngestionEndpoint=https://<location>.in.applicationinsights.azure.com/;LiveEndpoint=https://<location>.livediagnostics.monitor.azure.com/` (sensitive).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
