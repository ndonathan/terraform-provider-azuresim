---
subcategory: "Network"
page_title: "AzureSim: azuresim_virtual_network"
description: |-
  Manages a simulated Azure Virtual Network.
---

# azuresim_virtual_network

Manages a simulated Azure Virtual Network.

This resource mimics the behavior of the [`azurerm_virtual_network`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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

  dns_servers = [
    "10.0.0.4",
    "10.0.0.5",
  ]

  tags = {
    environment = "production"
  }
}
```

## Example Usage - Multiple Address Spaces

```terraform
resource "azuresim_virtual_network" "multi" {
  name                = "vnet-multi"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  address_space = [
    "10.0.0.0/16",
    "172.16.0.0/12",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Network. Changing this forces a new Virtual Network to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Virtual Network should exist. Changing this forces a new Virtual Network to be created.

* `location` - (Required) The Azure Region where the Virtual Network should exist. Changing this forces a new Virtual Network to be created.

* `address_space` - (Required) The list of address spaces used by the Virtual Network. Each element should be a valid CIDR block (e.g. `10.0.0.0/16`).

* `dns_servers` - (Optional) A list of IP addresses of DNS servers for the Virtual Network.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Network.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Network. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/virtualNetworks/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
