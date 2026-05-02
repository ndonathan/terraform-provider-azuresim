---
subcategory: "Network"
page_title: "AzureSim: azuresim_public_ip"
description: |-
  Manages a simulated Azure Public IP Address.
---

# azuresim_public_ip

Manages a simulated Azure Public IP Address.

This resource mimics the behavior of the [`azurerm_public_ip`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/public_ip) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_public_ip" "example" {
  name                = "pip-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"

  tags = {
    environment = "production"
  }
}
```

## Example Usage - With Domain Name Label

```terraform
resource "azuresim_public_ip" "labeled" {
  name                = "pip-labeled"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "myapp"
}

# fqdn => "myapp.eastus.cloudapp.azure.com"
output "fqdn" {
  value = azuresim_public_ip.labeled.fqdn
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Public IP. Changing this forces a new Public IP to be created.

* `resource_group_name` - (Required) The name of the Resource Group. Changing this forces a new Public IP to be created.

* `location` - (Required) The Azure Region. Changing this forces a new Public IP to be created.

* `allocation_method` - (Required) Defines the allocation method. Possible values are `Static` and `Dynamic`.

* `sku` - (Optional) The SKU of the Public IP. Possible values are `Basic` and `Standard`. Defaults to `Basic`.

* `sku_tier` - (Optional) The SKU tier. Possible values are `Regional` and `Global`. Defaults to `Regional`.

* `ip_version` - (Optional) The IP version. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`.

* `domain_name_label` - (Optional) A DNS label. When set, the resource's `fqdn` is computed as `<label>.<location>.cloudapp.azure.com`.

* `idle_timeout_in_minutes` - (Optional) The idle timeout in minutes. Must be between `4` and `30`.

* `reverse_fqdn` - (Optional) A fully qualified domain name reverse mapping.

* `zones` - (Optional) A list of availability zones (e.g. `["1", "2", "3"]`).

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Public IP. Follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/publicIPAddresses/{name}
  ```

* `ip_address` - A simulated IP address. For IPv4, an address in the `203.0.113.0/24` (TEST-NET-3) range deterministically derived from the resource name and resource group. For IPv6, an address in the `2001:db8::/32` documentation range.

* `fqdn` - The fully qualified domain name. Empty unless `domain_name_label` is set.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
