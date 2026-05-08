---
subcategory: "Database"
page_title: "AzureSim: azuresim_mssql_server"
description: |-
  Manages a simulated Azure SQL Server.
---

# azuresim_mssql_server

Manages a simulated Azure SQL Server.

This resource mimics the behavior of the [`azurerm_mssql_server`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_server) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_mssql_server" "example" {
  name                         = "sql-example-001"
  resource_group_name          = azuresim_resource_group.example.name
  location                     = azuresim_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "P@ssw0rd1234!"
  minimum_tls_version          = "1.2"
  public_network_access_enabled = true

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Server name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `version` - (Required) `2.0` (v11) or `12.0` (v12).
* `administrator_login` - (Optional) Server admin login (required if not using AAD-only auth).
* `administrator_login_password` - (Optional) Server admin password (sensitive).
* `minimum_tls_version` - (Optional) `1.0`, `1.1`, `1.2`, or `1.3`.
* `public_network_access_enabled` - (Optional) Whether public network access is enabled.
* `outbound_network_restriction_enabled` - (Optional) Restrict outbound network traffic.
* `tags` - (Optional) Tags.

~> **Note:** The `administrator_login_password` is stored in the Terraform state as sensitive data.

## Attributes Reference

* `id` - The SQL Server ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Sql/servers/{name}
  ```

* `fully_qualified_domain_name` - Simulated FQDN: `<name>.database.windows.net`.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
