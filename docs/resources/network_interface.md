---
subcategory: "Network"
page_title: "AzureSim: azuresim_network_interface"
description: |-
  Manages a simulated Azure Network Interface.
---

# azuresim_network_interface

Manages a simulated Azure Network Interface (NIC).

This resource mimics the behavior of the [`azurerm_network_interface`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_interface) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

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
  name                 = "snet-example"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azuresim_network_interface" "example" {
  name                = "nic-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azuresim_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "production"
  }
}
```

## Example Usage - Static IP with Public IP attached

```terraform
resource "azuresim_public_ip" "example" {
  name                = "pip-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azuresim_network_interface" "static" {
  name                          = "nic-static"
  resource_group_name           = azuresim_resource_group.example.name
  location                      = azuresim_resource_group.example.location
  accelerated_networking_enabled = true

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azuresim_subnet.example.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.1.10"
    public_ip_address_id          = azuresim_public_ip.example.id
    primary                       = true
  }

  dns_servers = ["10.0.0.4", "10.0.0.5"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Network Interface. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the Resource Group. Changing this forces a new resource to be created.

* `location` - (Required) Azure Region. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as defined below.

* `dns_servers` - (Optional) List of DNS server IP addresses.

* `internal_dns_name_label` - (Optional) Relative DNS name for this NIC inside the VNet.

* `accelerated_networking_enabled` - (Optional) Whether accelerated networking is enabled. Defaults to `false`.

* `ip_forwarding_enabled` - (Optional) Whether IP forwarding is enabled. Defaults to `false`.

* `edge_zone` - (Optional) Edge Zone within the Azure region. Changing this forces a new resource to be created.

* `tags` - (Optional) Mapping of tags.

---

An `ip_configuration` block supports:

* `name` - (Required) Name of the IP configuration.

* `subnet_id` - (Optional) Subnet ID. Required for IPv4 configurations.

* `private_ip_address_allocation` - (Required) `Static` or `Dynamic`.

* `private_ip_address` - (Optional) The private IP. Required when `private_ip_address_allocation` is `Static`. Computed when `Dynamic`.

* `private_ip_address_version` - (Optional) `IPv4` or `IPv6`. Defaults to `IPv4`.

* `public_ip_address_id` - (Optional) ID of an associated Public IP.

* `primary` - (Optional) Whether this is the primary IP configuration. Defaults to `true` for the first block, `false` thereafter.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The Network Interface ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/networkInterfaces/{name}
  ```

* `mac_address` - Simulated MAC address (Hyper-V OUI prefix `00-15-5D-…`, deterministic from name + RG).

* `private_ip_address` - Primary private IP, mirrored from the primary `ip_configuration`.

* `private_ip_addresses` - All private IPs, in `ip_configuration` block order.

* `applied_dns_servers` - Echoes `dns_servers`.

* `internal_domain_name_suffix` - Always `internal.cloudapp.net`.

When `private_ip_address_allocation` is `Dynamic`, simulated IPs are deterministically derived from the NIC name + IP-config name and placed in `10.0.0.0/24` (IPv4) or `fd00::/8` (IPv6). They do **not** validate against the subnet's actual address prefix — pass `private_ip_address` explicitly if you need a specific address.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
