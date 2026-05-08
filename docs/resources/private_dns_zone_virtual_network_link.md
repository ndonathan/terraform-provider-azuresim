---
subcategory: "Network"
page_title: "AzureSim: azuresim_private_dns_zone_virtual_network_link"
description: |-
  Links a Virtual Network to a simulated Private DNS Zone.
---

# azuresim_private_dns_zone_virtual_network_link

Links a Virtual Network to a simulated Azure Private DNS Zone.

This resource mimics the behavior of the [`azurerm_private_dns_zone_virtual_network_link`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_dns_zone_virtual_network_link) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_virtual_network" "example" {
  name                = "vnet-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azuresim_private_dns_zone" "example" {
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = azuresim_resource_group.example.name
}

resource "azuresim_private_dns_zone_virtual_network_link" "example" {
  name                  = "vnet-link"
  resource_group_name   = azuresim_resource_group.example.name
  private_dns_zone_name = azuresim_private_dns_zone.example.name
  virtual_network_id    = azuresim_virtual_network.example.id
  registration_enabled  = false

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Link name.
* `resource_group_name` - (Required, ForceNew) Resource Group of the parent Private DNS Zone.
* `private_dns_zone_name` - (Required, ForceNew) Parent Private DNS Zone name.
* `virtual_network_id` - (Required, ForceNew) Virtual Network ID.
* `registration_enabled` - (Optional) Enable auto-registration of VM hostnames in the zone.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Link ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/privateDnsZones/{zone_name}/virtualNetworkLinks/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
