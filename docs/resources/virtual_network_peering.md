---
subcategory: "Network"
page_title: "AzureSim: azuresim_virtual_network_peering"
description: |-
  Manages a simulated Virtual Network Peering.
---

# azuresim_virtual_network_peering

Manages a simulated Virtual Network Peering. Create one resource per direction (peering is unidirectional in the data model).

This resource mimics the behavior of the [`azurerm_virtual_network_peering`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network_peering) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_virtual_network" "hub" {
  name                = "vnet-hub"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azuresim_virtual_network" "spoke" {
  name                = "vnet-spoke"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  address_space       = ["10.1.0.0/16"]
}

resource "azuresim_virtual_network_peering" "hub_to_spoke" {
  name                         = "hub-to-spoke"
  resource_group_name          = azuresim_resource_group.example.name
  virtual_network_name         = azuresim_virtual_network.hub.name
  remote_virtual_network_id    = azuresim_virtual_network.spoke.id
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true
}

resource "azuresim_virtual_network_peering" "spoke_to_hub" {
  name                         = "spoke-to-hub"
  resource_group_name          = azuresim_resource_group.example.name
  virtual_network_name         = azuresim_virtual_network.spoke.name
  remote_virtual_network_id    = azuresim_virtual_network.hub.id
  allow_virtual_network_access = true
  use_remote_gateways          = true
}
```

## Argument Reference

* `name` - (Required, ForceNew) Peering name.
* `resource_group_name` - (Required, ForceNew) Resource Group of the local VNet.
* `virtual_network_name` - (Required, ForceNew) Local VNet name.
* `remote_virtual_network_id` - (Required, ForceNew) Remote VNet ID.
* `allow_virtual_network_access` - (Optional) Allow access to the remote VNet.
* `allow_forwarded_traffic` - (Optional) Allow forwarded traffic from the remote VNet.
* `allow_gateway_transit` - (Optional) Allow gateway transit.
* `use_remote_gateways` - (Optional) Use the remote VNet's gateway.

## Attributes Reference

* `id` - The Peering ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/virtualNetworks/{vnet_name}/virtualNetworkPeerings/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
