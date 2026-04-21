---
subcategory: "Network"
page_title: "AzureSim: azuresim_subnet"
description: |-
  Manages a simulated Azure Subnet.
---

# azuresim_subnet

Manages a simulated Azure Subnet within a Virtual Network.

This resource mimics the behavior of the [`azurerm_subnet`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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

resource "azuresim_subnet" "example" {
  name                 = "snet-internal"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}
```

## Example Usage - Multiple Subnets

```terraform
resource "azuresim_subnet" "frontend" {
  name                 = "snet-frontend"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azuresim_subnet" "backend" {
  name                 = "snet-backend"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Subnet. Changing this forces a new Subnet to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Subnet should exist. Changing this forces a new Subnet to be created.

* `virtual_network_name` - (Required) The name of the Virtual Network in which the Subnet should exist. Changing this forces a new Subnet to be created.

* `address_prefixes` - (Required) The address prefixes to use for the Subnet. Each element should be a valid CIDR block (e.g. `10.0.1.0/24`).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Subnet. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/virtualNetworks/{virtual_network_name}/subnets/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
