---
subcategory: "Database"
page_title: "AzureSim: azuresim_postgresql_flexible_server"
description: |-
  Manages a simulated Azure PostgreSQL Flexible Server.
---

# azuresim_postgresql_flexible_server

Manages a simulated Azure PostgreSQL Flexible Server.

This resource mimics the behavior of the [`azurerm_postgresql_flexible_server`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_flexible_server) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_postgresql_flexible_server" "example" {
  name                   = "psql-example"
  resource_group_name    = azuresim_resource_group.example.name
  location               = azuresim_resource_group.example.location
  version                = "16"
  sku_name               = "GP_Standard_D2s_v3"
  storage_mb             = 32768
  backup_retention_days  = 14
  administrator_login    = "psqladmin"
  administrator_password = "P@ssw0rd1234!"
  zone                   = "1"

  public_network_access_enabled = true

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Server name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `version` - (Optional) Postgres major version (e.g. `14`, `15`, `16`).
* `sku_name` - (Required) SKU (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`).
* `storage_mb` - (Optional) Storage in MB.
* `storage_tier` - (Optional) Storage tier (e.g. `P4`, `P6`, `P10`).
* `backup_retention_days` - (Optional) Backup retention (7-35).
* `geo_redundant_backup_enabled` - (Optional) Enable geo-redundant backup.
* `administrator_login` - (Optional) Admin login.
* `administrator_password` - (Optional) Admin password (sensitive).
* `delegated_subnet_id` - (Optional) Delegated subnet ID for VNet integration.
* `private_dns_zone_id` - (Optional) Private DNS zone ID.
* `zone` - (Optional) Availability zone.
* `public_network_access_enabled` - (Optional) Allow public access.
* `high_availability` - (Optional) One `high_availability` block as defined below.
* `tags` - (Optional) Tags.

---

A `high_availability` block supports:

* `mode` - (Required) `SameZone` or `ZoneRedundant`.
* `standby_availability_zone` - (Optional) Standby zone.

## Attributes Reference

* `id` - The PostgreSQL Flexible Server ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{name}
  ```

* `fqdn` - Simulated FQDN: `<name>.postgres.database.azure.com`.

~> **Note:** The `administrator_password` is stored in the Terraform state as sensitive data.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
