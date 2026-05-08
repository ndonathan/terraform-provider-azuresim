---
subcategory: "Network"
page_title: "AzureSim: azuresim_subnet_route_table_association"
description: |-
  Associates a simulated Route Table with a Subnet.
---

# azuresim_subnet_route_table_association

Associates a simulated Route Table with a Subnet.

This resource mimics the behavior of the [`azurerm_subnet_route_table_association`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet_route_table_association) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

Mirrors the AzureRM convention of using the Subnet ID as the resource ID.

## Example Usage

```terraform
resource "azuresim_subnet" "example" {
  name                 = "snet-example"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azuresim_route_table" "example" {
  name                = "rt-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
}

resource "azuresim_subnet_route_table_association" "example" {
  subnet_id      = azuresim_subnet.example.id
  route_table_id = azuresim_route_table.example.id
}
```

## Argument Reference

* `subnet_id` - (Required, ForceNew) Subnet ID.
* `route_table_id` - (Required) Route Table ID.

## Attributes Reference

* `id` - Equal to the Subnet ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/virtualNetworks/{vnet_name}/subnets/{subnet_name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
