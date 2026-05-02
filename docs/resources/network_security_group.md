---
subcategory: "Network"
page_title: "AzureSim: azuresim_network_security_group"
description: |-
  Manages a simulated Azure Network Security Group.
---

# azuresim_network_security_group

Manages a simulated Azure Network Security Group.

This resource mimics the behavior of the [`azurerm_network_security_group`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_security_group) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_network_security_group" "example" {
  name                = "nsg-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  security_rule {
    name                       = "AllowSSH"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowHTTPS"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "Internet"
    destination_address_prefix = "VirtualNetwork"
  }

  tags = {
    environment = "production"
  }
}
```

## Example Usage - Multiple Ports and Prefixes

```terraform
resource "azuresim_network_security_group" "multi" {
  name                = "nsg-multi"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  security_rule {
    name                         = "AllowWeb"
    priority                     = 200
    direction                    = "Inbound"
    access                       = "Allow"
    protocol                     = "Tcp"
    source_port_range            = "*"
    destination_port_ranges      = ["80", "443", "8080"]
    source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
    destination_address_prefix   = "*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Security Group. Changing this forces a new Network Security Group to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Network Security Group should exist. Changing this forces a new Network Security Group to be created.

* `location` - (Required) The Azure Region where the Network Security Group should exist. Changing this forces a new Network Security Group to be created.

* `security_rule` - (Optional) One or more `security_rule` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Security Group.

---

A `security_rule` block supports the following:

* `name` - (Required) The name of the security rule.

* `description` - (Optional) A description for the rule.

* `protocol` - (Required) The network protocol this rule applies to. Possible values are `Tcp`, `Udp`, `Icmp`, `Esp`, `Ah`, or `*`.

* `source_port_range` - (Optional) Source port or range. Use `*` to match any port. Mutually exclusive with `source_port_ranges`.

* `source_port_ranges` - (Optional) A list of source ports or port ranges. Mutually exclusive with `source_port_range`.

* `destination_port_range` - (Optional) Destination port or range. Mutually exclusive with `destination_port_ranges`.

* `destination_port_ranges` - (Optional) A list of destination ports or port ranges. Mutually exclusive with `destination_port_range`.

* `source_address_prefix` - (Optional) CIDR or service tag (e.g. `*`, `VirtualNetwork`, `Internet`). Mutually exclusive with `source_address_prefixes`.

* `source_address_prefixes` - (Optional) A list of source CIDRs. Mutually exclusive with `source_address_prefix`.

* `destination_address_prefix` - (Optional) CIDR or service tag. Mutually exclusive with `destination_address_prefixes`.

* `destination_address_prefixes` - (Optional) A list of destination CIDRs. Mutually exclusive with `destination_address_prefix`.

* `access` - (Required) Whether traffic is allowed or denied. Possible values are `Allow` and `Deny`.

* `priority` - (Required) Rule priority between `100` and `4096`. Lower numbers are evaluated first.

* `direction` - (Required) Direction of the rule. Possible values are `Inbound` and `Outbound`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Security Group. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/networkSecurityGroups/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
