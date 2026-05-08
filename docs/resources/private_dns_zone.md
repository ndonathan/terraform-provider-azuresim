---
subcategory: "Network"
page_title: "AzureSim: azuresim_private_dns_zone"
description: |-
  Manages a simulated Azure Private DNS Zone.
---

# azuresim_private_dns_zone

Manages a simulated Azure Private DNS Zone.

This resource mimics the behavior of the [`azurerm_private_dns_zone`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_dns_zone) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_private_dns_zone" "example" {
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = azuresim_resource_group.example.name

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Zone name (e.g. `privatelink.blob.core.windows.net`).
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Private DNS Zone ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/privateDnsZones/{name}
  ```

* `number_of_record_sets` - Always `0` in this simulator.
* `max_number_of_record_sets` - Static (`25000`).
* `number_of_virtual_network_links` - Always `0` in this simulator.
* `max_number_of_virtual_network_links` - Static (`1000`).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
