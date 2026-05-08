---
subcategory: "Network"
page_title: "AzureSim: azuresim_bastion_host"
description: |-
  Manages a simulated Azure Bastion Host.
---

# azuresim_bastion_host

Manages a simulated Azure Bastion Host.

This resource mimics the behavior of the [`azurerm_bastion_host`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/bastion_host) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_subnet" "bastion" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.255.0/27"]
}

resource "azuresim_public_ip" "bastion" {
  name                = "pip-bastion"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azuresim_bastion_host" "example" {
  name                = "bst-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azuresim_subnet.bastion.id
    public_ip_address_id = azuresim_public_ip.bastion.id
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the Bastion Host.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Optional) SKU. One of `Basic`, `Standard`, `Developer`, or `Premium`.
* `scale_units` - (Optional) Scale units for `Standard` SKU (2-50).
* `copy_paste_enabled` - (Optional) Whether copy/paste is allowed.
* `file_copy_enabled` - (Optional) Whether file copy is allowed (`Standard` SKU only).
* `ip_connect_enabled` - (Optional) Whether IP-based connection is allowed.
* `shareable_link_enabled` - (Optional) Whether shareable links are allowed.
* `tunneling_enabled` - (Optional) Whether native client tunneling is allowed.
* `ip_configuration` - (Required) One `ip_configuration` block as defined below. The associated subnet must be named `AzureBastionSubnet`.
* `tags` - (Optional) Tags.

---

An `ip_configuration` block supports:

* `name` - (Required) Configuration name.
* `subnet_id` - (Required) ID of the `AzureBastionSubnet`.
* `public_ip_address_id` - (Required) Public IP address ID.

## Attributes Reference

* `id` - The Bastion Host ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/bastionHosts/{name}
  ```

* `dns_name` - Simulated DNS name of the form `bst-<hash>.bastion.azure.com`, deterministically derived from the host name.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
