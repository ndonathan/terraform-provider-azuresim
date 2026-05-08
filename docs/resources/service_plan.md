---
subcategory: "Web"
page_title: "AzureSim: azuresim_service_plan"
description: |-
  Manages a simulated Azure App Service Plan.
---

# azuresim_service_plan

Manages a simulated Azure App Service Plan (`Microsoft.Web/serverfarms`).

This resource mimics the behavior of the [`azurerm_service_plan`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_plan) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_service_plan" "example" {
  name                = "asp-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "P1v3"
  worker_count        = 2
  zone_balancing_enabled = true

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the Service Plan.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `os_type` - (Required, ForceNew) `Linux`, `Windows`, or `WindowsContainer`.
* `sku_name` - (Required) SKU (e.g. `B1`, `S1`, `P1v3`, `Y1` for Consumption Functions, `EP1` for Elastic Premium).
* `worker_count` - (Optional) Number of workers.
* `maximum_elastic_worker_count` - (Optional) Max elastic worker count (Premium plans only).
* `per_site_scaling_enabled` - (Optional) Enable per-site scaling.
* `zone_balancing_enabled` - (Optional) Enable zone balancing.
* `app_service_environment_id` - (Optional) ASE ID for ASEv3 plans.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Service Plan ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Web/serverfarms/{name}
  ```

* `reserved` - Whether the plan is Linux (`true` for `os_type = "Linux"`, otherwise `false`).
* `kind` - Plan kind: `linux` (for Linux), `windows,container` (for `WindowsContainer`), or `app` (default Windows).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
