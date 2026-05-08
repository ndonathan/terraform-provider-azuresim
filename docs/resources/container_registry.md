---
subcategory: "Container"
page_title: "AzureSim: azuresim_container_registry"
description: |-
  Manages a simulated Azure Container Registry.
---

# azuresim_container_registry

Manages a simulated Azure Container Registry (ACR).

This resource mimics the behavior of the [`azurerm_container_registry`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/container_registry) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_container_registry" "example" {
  name                = "acrexample001"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Premium"
  admin_enabled       = true

  zone_redundancy_enabled       = true
  public_network_access_enabled = true

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Registry name. In a real Azure environment, must be globally unique, 5-50 alphanumeric characters.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Required) `Basic`, `Standard`, or `Premium`.
* `admin_enabled` - (Optional) Enable the admin user. Defaults to `false`.
* `public_network_access_enabled` - (Optional) Allow public network access. Defaults to `true`.
* `zone_redundancy_enabled` - (Optional) Enable zone redundancy (`Premium` only).
* `anonymous_pull_enabled` - (Optional) Allow anonymous pulls.
* `data_endpoint_enabled` - (Optional) Enable dedicated data endpoints.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Container Registry ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ContainerRegistry/registries/{name}
  ```

* `login_server` - Simulated login server (`<name>.azurecr.io`).
* `admin_username` - Admin username (always equals `name` when admin is enabled).
* `admin_password` - Simulated admin password (sensitive, deterministically derived from RG + name).

~> **Note:** The `admin_password` is stored in the Terraform state as sensitive data. The value is a static placeholder and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
