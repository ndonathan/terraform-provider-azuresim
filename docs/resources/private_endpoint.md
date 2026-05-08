---
subcategory: "Network"
page_title: "AzureSim: azuresim_private_endpoint"
description: |-
  Manages a simulated Azure Private Endpoint.
---

# azuresim_private_endpoint

Manages a simulated Azure Private Endpoint.

This resource mimics the behavior of the [`azurerm_private_endpoint`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_storage_account" "example" {
  name                     = "stpeexample"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_private_endpoint" "example" {
  name                = "pe-storage"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  subnet_id           = azuresim_subnet.example.id

  private_service_connection {
    name                           = "psc-blob"
    private_connection_resource_id = azuresim_storage_account.example.id
    is_manual_connection           = false
    subresource_names              = ["blob"]
  }

  private_dns_zone_group {
    name                 = "default"
    private_dns_zone_ids = [azuresim_private_dns_zone.example.id]
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Endpoint name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `subnet_id` - (Required, ForceNew) Subnet ID where the endpoint NIC is placed.
* `custom_network_interface_name` - (Optional, ForceNew) Custom name for the auto-created NIC.
* `private_service_connection` - (Required) One or more `private_service_connection` blocks as defined below.
* `private_dns_zone_group` - (Optional) Zero or more `private_dns_zone_group` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `private_service_connection` block supports:

* `name` - (Required) Connection name.
* `private_connection_resource_id` - (Optional) Target resource ID. Mutually exclusive with `private_connection_resource_alias`.
* `private_connection_resource_alias` - (Optional) Target resource alias.
* `is_manual_connection` - (Required) Whether the connection requires manual approval.
* `subresource_names` - (Optional) Subresource names (e.g. `["blob"]`, `["sqlServer"]`).
* `request_message` - (Optional) Approval request message (manual approval only).

---

A `private_dns_zone_group` block supports:

* `name` - (Required) Group name.
* `private_dns_zone_ids` - (Required) Private DNS zone IDs.

## Attributes Reference

* `id` - The Private Endpoint ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/privateEndpoints/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
