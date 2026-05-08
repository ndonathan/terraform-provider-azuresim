---
subcategory: "Web"
page_title: "AzureSim: azuresim_windows_function_app"
description: |-
  Manages a simulated Azure Windows Function App.
---

# azuresim_windows_function_app

Manages a simulated Azure Windows Function App.

This resource mimics the behavior of the [`azurerm_windows_function_app`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_function_app) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_storage_account" "example" {
  name                     = "stfunctionexample"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_service_plan" "example" {
  name                = "asp-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  os_type             = "Windows"
  sku_name            = "Y1"
}

resource "azuresim_windows_function_app" "example" {
  name                       = "func-example"
  resource_group_name        = azuresim_resource_group.example.name
  location                   = azuresim_resource_group.example.location
  service_plan_id            = azuresim_service_plan.example.id
  storage_account_name       = azuresim_storage_account.example.name
  storage_account_access_key = azuresim_storage_account.example.primary_access_key

  https_only                  = true
  functions_extension_version = "~4"

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "dotnet-isolated"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Function App name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `service_plan_id` - (Required) Parent App Service Plan ID.
* `storage_account_name` - (Required) Backing Storage Account name.
* `storage_account_access_key` - (Optional) Backing Storage Account access key. Mutually exclusive with `storage_uses_managed_identity = true`.
* `storage_uses_managed_identity` - (Optional) Use the system-assigned identity to access the storage account.
* `https_only` - (Optional) Force HTTPS.
* `functions_extension_version` - (Optional) Functions runtime version (e.g. `~4`).
* `builtin_logging_enabled` - (Optional) Enable Application Insights built-in logging.
* `app_settings` - (Optional) Map of environment variables.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Function App ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Web/sites/{name}
  ```

* `default_hostname` - Simulated default hostname (`<name>.azurewebsites.net`).
* `outbound_ip_addresses` - Simulated outbound IPs (comma-separated).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
