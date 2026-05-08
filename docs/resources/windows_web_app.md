---
subcategory: "Web"
page_title: "AzureSim: azuresim_windows_web_app"
description: |-
  Manages a simulated Azure Windows Web App.
---

# azuresim_windows_web_app

Manages a simulated Azure Windows Web App.

This resource mimics the behavior of the [`azurerm_windows_web_app`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_web_app) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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
  os_type             = "Windows"
  sku_name            = "P1v3"
}

resource "azuresim_windows_web_app" "example" {
  name                    = "app-example"
  resource_group_name     = azuresim_resource_group.example.name
  location                = azuresim_resource_group.example.location
  service_plan_id         = azuresim_service_plan.example.id
  https_only              = true

  site_config {
    always_on           = true
    http2_enabled       = true
    minimum_tls_version = "1.2"
    windows_fx_version  = "DOTNETCORE|8.0"
  }

  app_settings = {
    ASPNETCORE_ENVIRONMENT = "Production"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Web App name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `service_plan_id` - (Required) Parent App Service Plan ID.
* `https_only` - (Optional) Force HTTPS.
* `client_affinity_enabled` - (Optional) Enable session affinity.
* `app_settings` - (Optional) Map of environment variables.
* `site_config` - (Optional) One `site_config` block as defined below.
* `tags` - (Optional) Tags.

---

A `site_config` block supports:

* `always_on` - (Optional) Keep the app always on.
* `ftps_state` - (Optional) `AllAllowed`, `FtpsOnly`, or `Disabled`.
* `http2_enabled` - (Optional) Enable HTTP/2.
* `minimum_tls_version` - (Optional) `1.0`, `1.1`, `1.2`, or `1.3`.
* `websockets_enabled` - (Optional) Enable WebSockets.
* `health_check_path` - (Optional) Health check path.
* `vnet_route_all_enabled` - (Optional) Route all outbound traffic through the VNet.
* `linux_fx_version` - (Optional) Linux runtime stack (typically not used here).
* `windows_fx_version` - (Optional) Windows runtime stack (e.g. `DOTNETCORE|8.0`).
* `app_command_line` - (Optional) Startup command.

## Attributes Reference

* `id` - The Web App ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Web/sites/{name}
  ```

* `default_hostname` - Simulated default hostname (`<name>.azurewebsites.net`).
* `outbound_ip_addresses` - Simulated outbound IPs (comma-separated).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
