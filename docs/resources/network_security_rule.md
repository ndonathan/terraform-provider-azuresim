---
subcategory: "Network"
page_title: "AzureSim: azuresim_network_security_rule"
description: |-
  Manages a standalone simulated NSG security rule.
---

# azuresim_network_security_rule

Manages a standalone simulated Network Security Group rule.

This resource mimics the behavior of the [`azurerm_network_security_rule`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_security_rule) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

~> **Note:** Use either inline `security_rule` blocks on `azuresim_network_security_group` *or* this resource — not both.

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
}

resource "azuresim_network_security_rule" "allow_ssh" {
  name                        = "allow-ssh"
  resource_group_name         = azuresim_resource_group.example.name
  network_security_group_name = azuresim_network_security_group.example.name

  protocol                   = "Tcp"
  access                     = "Allow"
  direction                  = "Inbound"
  priority                   = 100
  source_port_range          = "*"
  destination_port_range     = "22"
  source_address_prefix      = "*"
  destination_address_prefix = "VirtualNetwork"
}
```

## Argument Reference

* `name` - (Required, ForceNew) Rule name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `network_security_group_name` - (Required, ForceNew) Parent NSG name.
* `description` - (Optional) Description.
* `protocol` - (Required) `Tcp`, `Udp`, `Icmp`, `Esp`, `Ah`, or `*`.
* `source_port_range` - (Optional) Source port or range. Mutually exclusive with `source_port_ranges`.
* `source_port_ranges` - (Optional) List of source ports.
* `destination_port_range` - (Optional) Destination port or range.
* `destination_port_ranges` - (Optional) List of destination ports.
* `source_address_prefix` - (Optional) Source CIDR or service tag (e.g. `Internet`, `VirtualNetwork`).
* `source_address_prefixes` - (Optional) List of source CIDRs.
* `destination_address_prefix` - (Optional) Destination CIDR or service tag.
* `destination_address_prefixes` - (Optional) List of destination CIDRs.
* `access` - (Required) `Allow` or `Deny`.
* `priority` - (Required) Priority (100-4096).
* `direction` - (Required) `Inbound` or `Outbound`.

## Attributes Reference

* `id` - The Security Rule ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/networkSecurityGroups/{nsg_name}/securityRules/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
