---
subcategory: "Network"
page_title: "AzureSim: azuresim_cdn_frontdoor_profile"
description: |-
  Manages a simulated Azure Front Door (Standard/Premium) profile.
---

# azuresim_cdn_frontdoor_profile

Manages a simulated Azure Front Door (Standard/Premium) profile.

This resource mimics the behavior of the [`azurerm_cdn_frontdoor_profile`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cdn_frontdoor_profile) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_cdn_frontdoor_profile" "example" {
  name                     = "afd-example"
  resource_group_name      = azuresim_resource_group.example.name
  sku_name                 = "Standard_AzureFrontDoor"
  response_timeout_seconds = 60

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Profile name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `sku_name` - (Required) `Standard_AzureFrontDoor` or `Premium_AzureFrontDoor`.
* `response_timeout_seconds` - (Optional) Origin response timeout (16-240).
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Front Door profile ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Cdn/profiles/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
