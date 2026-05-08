---
subcategory: "Network"
page_title: "AzureSim: azuresim_lb"
description: |-
  Manages a simulated Azure Load Balancer.
---

# azuresim_lb

Manages a simulated Azure Load Balancer.

This resource mimics the behavior of the [`azurerm_lb`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/lb) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage - Public Load Balancer

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_public_ip" "example" {
  name                = "pip-lb"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azuresim_lb" "example" {
  name                = "lb-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azuresim_public_ip.example.id
  }

  tags = {
    environment = "production"
  }
}
```

## Example Usage - Internal Load Balancer

```terraform
resource "azuresim_lb" "internal" {
  name                = "lb-internal"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard"

  frontend_ip_configuration {
    name                          = "internal"
    subnet_id                     = azuresim_subnet.example.id
    private_ip_address            = "10.0.1.100"
    private_ip_address_allocation = "Static"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Load Balancer name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Optional) `Basic`, `Standard`, or `Gateway`. Defaults to `Basic`.
* `sku_tier` - (Optional) `Regional` or `Global`.
* `edge_zone` - (Optional, ForceNew) Edge zone.
* `frontend_ip_configuration` - (Required) One or more `frontend_ip_configuration` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `frontend_ip_configuration` block supports:

* `name` - (Required) Configuration name.
* `public_ip_address_id` - (Optional) Public IP ID (for public LBs).
* `subnet_id` - (Optional) Subnet ID (for internal LBs).
* `private_ip_address` - (Optional) Static private IP (when allocation is `Static`).
* `private_ip_address_allocation` - (Optional) `Static` or `Dynamic`.
* `private_ip_address_version` - (Optional) `IPv4` or `IPv6`.
* `zones` - (Optional) Availability zones.

## Attributes Reference

* `id` - The Load Balancer ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/loadBalancers/{name}
  ```

* `private_ip_address` - Primary private IP, mirrored from the first frontend with a private allocation.
* `private_ip_addresses` - All private IPs across frontend configurations.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
