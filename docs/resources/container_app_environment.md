---
subcategory: "Container"
page_title: "AzureSim: azuresim_container_app_environment"
description: |-
  Manages a simulated Azure Container App Environment.
---

# azuresim_container_app_environment

Manages a simulated Azure Container App Environment.

This resource mimics the behavior of the [`azurerm_container_app_environment`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/container_app_environment) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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

resource "azuresim_container_app_environment" "example" {
  name                       = "cae-example"
  resource_group_name        = azuresim_resource_group.example.name
  location                   = azuresim_resource_group.example.location
  log_analytics_workspace_id = azuresim_log_analytics_workspace.example.id
  zone_redundancy_enabled    = true
  workload_profile_type      = "Consumption"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Environment name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `log_analytics_workspace_id` - (Optional) Linked Log Analytics Workspace ID.
* `infrastructure_subnet_id` - (Optional, ForceNew) Infrastructure subnet ID.
* `internal_load_balancer_enabled` - (Optional) Whether to use an internal load balancer.
* `zone_redundancy_enabled` - (Optional) Enable zone redundancy.
* `workload_profile_type` - (Optional) Workload profile (e.g. `Consumption`, `D4`).
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Container App Environment ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.App/managedEnvironments/{name}
  ```

* `default_domain` - Simulated default domain of the form `<name>.<location>.azurecontainerapps.io`.
* `static_ip_address` - Simulated static IP address (`203.0.113.50`).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
