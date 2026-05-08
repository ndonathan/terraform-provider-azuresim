---
subcategory: "Database"
page_title: "AzureSim: azuresim_mssql_database"
description: |-
  Manages a simulated Azure SQL Database.
---

# azuresim_mssql_database

Manages a simulated Azure SQL Database.

This resource mimics the behavior of the [`azurerm_mssql_database`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_database) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_mssql_server" "example" {
  name                         = "sql-example-001"
  resource_group_name          = azuresim_resource_group.example.name
  location                     = azuresim_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "P@ssw0rd1234!"
}

resource "azuresim_mssql_database" "example" {
  name           = "appdb"
  server_id      = azuresim_mssql_server.example.id
  sku_name       = "GP_Gen5_2"
  collation      = "SQL_Latin1_General_CP1_CI_AS"
  max_size_gb    = 32
  zone_redundant = false
  license_type   = "LicenseIncluded"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Database name.
* `server_id` - (Required, ForceNew) Parent SQL Server ID.
* `sku_name` - (Optional) SKU (e.g. `Basic`, `S0`, `P1`, `GP_Gen5_2`, `BC_Gen5_4`).
* `collation` - (Optional) Collation (e.g. `SQL_Latin1_General_CP1_CI_AS`).
* `max_size_gb` - (Optional) Max database size in GB.
* `zone_redundant` - (Optional) Enable zone redundancy.
* `geo_backup_enabled` - (Optional) Enable geo-redundant backups.
* `storage_account_type` - (Optional) `Geo`, `GeoZone`, `Local`, or `Zone`.
* `read_scale` - (Optional) Enable read-scale (Premium/Business Critical only).
* `license_type` - (Optional) `LicenseIncluded` or `BasePrice`.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The SQL Database ID, formed as `<server_id>/databases/<name>`. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Sql/servers/{server_name}/databases/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
