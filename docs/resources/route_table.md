---
subcategory: "Network"
page_title: "AzureSim: azuresim_route_table"
description: |-
  Manages a simulated Azure Route Table.
---

# azuresim_route_table

Manages a simulated Azure Route Table.

This resource mimics the behavior of the [`azurerm_route_table`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/route_table) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_route_table" "example" {
  name                          = "rt-example"
  resource_group_name           = azuresim_resource_group.example.name
  location                      = azuresim_resource_group.example.location
  bgp_route_propagation_enabled = false

  route {
    name           = "default"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VirtualAppliance"
    next_hop_in_ip_address = "10.0.0.4"
  }

  route {
    name           = "vnetlocal"
    address_prefix = "10.0.0.0/16"
    next_hop_type  = "VnetLocal"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Route Table name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `bgp_route_propagation_enabled` - (Optional) Whether BGP route propagation is enabled. Defaults to `true`.
* `route` - (Optional) Zero or more `route` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `route` block supports:

* `name` - (Required) Route name.
* `address_prefix` - (Required) Destination CIDR.
* `next_hop_type` - (Required) `VirtualNetworkGateway`, `VnetLocal`, `Internet`, `VirtualAppliance`, or `None`.
* `next_hop_in_ip_address` - (Optional) Required when `next_hop_type` is `VirtualAppliance`.

## Attributes Reference

* `id` - The Route Table ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/routeTables/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
