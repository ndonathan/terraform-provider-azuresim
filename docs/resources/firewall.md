---
subcategory: "Network"
page_title: "AzureSim: azuresim_firewall"
description: |-
  Manages a simulated Azure Firewall.
---

# azuresim_firewall

Manages a simulated Azure Firewall.

This resource mimics the behavior of the [`azurerm_firewall`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/firewall) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_subnet" "firewall" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.255.0/26"]
}

resource "azuresim_public_ip" "firewall" {
  name                = "pip-firewall"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azuresim_firewall" "example" {
  name                = "afw-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"
  threat_intel_mode   = "Alert"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azuresim_subnet.firewall.id
    public_ip_address_id = azuresim_public_ip.firewall.id
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Firewall name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku_name` - (Required) `AZFW_VNet` or `AZFW_Hub`.
* `sku_tier` - (Required) `Standard`, `Premium`, or `Basic`.
* `firewall_policy_id` - (Optional) Firewall policy ID.
* `threat_intel_mode` - (Optional) `Off`, `Alert`, or `Deny`.
* `dns_servers` - (Optional) Custom DNS servers.
* `private_ip_ranges` - (Optional) Private IP ranges (SNAT exemptions).
* `zones` - (Optional) Availability zones.
* `ip_configuration` - (Optional) One or more `ip_configuration` blocks. The associated subnet must be named `AzureFirewallSubnet`.
* `management_ip_configuration` - (Optional) One `management_ip_configuration` block (Forced Tunneling). Subnet must be named `AzureFirewallManagementSubnet`.
* `tags` - (Optional) Tags.

---

An `ip_configuration` or `management_ip_configuration` block supports:

* `name` - (Required) Configuration name.
* `subnet_id` - (Optional) Subnet ID (`AzureFirewallSubnet` for `ip_configuration`).
* `public_ip_address_id` - (Optional) Public IP ID.

## Attributes Reference

* `id` - The Firewall ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/azureFirewalls/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
