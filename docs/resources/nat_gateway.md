---
subcategory: "Network"
page_title: "AzureSim: azuresim_nat_gateway"
description: |-
  Manages a simulated Azure NAT Gateway.
---

# azuresim_nat_gateway

Manages a simulated Azure NAT Gateway.

This resource mimics the behavior of the [`azurerm_nat_gateway`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/nat_gateway) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_nat_gateway" "example" {
  name                    = "natgw-example"
  resource_group_name     = azuresim_resource_group.example.name
  location                = azuresim_resource_group.example.location
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) NAT Gateway name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku_name` - (Optional) Always `Standard`.
* `idle_timeout_in_minutes` - (Optional) Idle timeout (4-120 minutes).
* `zones` - (Optional) Availability zones.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The NAT Gateway ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/natGateways/{name}
  ```

* `resource_guid` - Simulated resource GUID, deterministically derived from RG + name.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
