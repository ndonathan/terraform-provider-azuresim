---
subcategory: "Network"
page_title: "AzureSim: azuresim_subnet_network_security_group_association"
description: |-
  Associates a simulated Network Security Group with a Subnet.
---

# azuresim_subnet_network_security_group_association

Associates a simulated Network Security Group with a Subnet.

This resource mimics the behavior of the [`azurerm_subnet_network_security_group_association`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet_network_security_group_association) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

Mirrors the AzureRM convention of using the Subnet ID as the resource ID.

## Example Usage

```terraform
resource "azuresim_subnet" "example" {
  name                 = "snet-example"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azuresim_network_security_group" "example" {
  name                = "nsg-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
}

resource "azuresim_subnet_network_security_group_association" "example" {
  subnet_id                 = azuresim_subnet.example.id
  network_security_group_id = azuresim_network_security_group.example.id
}
```

## Argument Reference

* `subnet_id` - (Required, ForceNew) Subnet ID.
* `network_security_group_id` - (Required) Network Security Group ID.

## Attributes Reference

* `id` - Equal to the Subnet ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/virtualNetworks/{vnet_name}/subnets/{subnet_name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
