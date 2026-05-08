---
subcategory: "Database"
page_title: "AzureSim: azuresim_mysql_flexible_server"
description: |-
  Manages a simulated Azure MySQL Flexible Server.
---

# azuresim_mysql_flexible_server

Manages a simulated Azure MySQL Flexible Server.

This resource mimics the behavior of the [`azurerm_mysql_flexible_server`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mysql_flexible_server) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_mysql_flexible_server" "example" {
  name                   = "mysql-example"
  resource_group_name    = azuresim_resource_group.example.name
  location               = azuresim_resource_group.example.location
  version                = "8.0.21"
  sku_name               = "GP_Standard_D2ds_v4"
  backup_retention_days  = 14
  administrator_login    = "mysqladmin"
  administrator_password = "P@ssw0rd1234!"
  zone                   = "1"

  storage {
    size_gb           = 32
    auto_grow_enabled = true
  }

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
* `version` - (Optional) MySQL version (e.g. `5.7`, `8.0.21`).
* `sku_name` - (Required) SKU.
* `backup_retention_days` - (Optional) Backup retention (1-35).
* `geo_redundant_backup_enabled` - (Optional) Enable geo-redundant backup.
* `administrator_login` - (Optional) Admin login.
* `administrator_password` - (Optional) Admin password (sensitive).
* `delegated_subnet_id` - (Optional) Delegated subnet ID for VNet integration.
* `private_dns_zone_id` - (Optional) Private DNS zone ID.
* `zone` - (Optional) Availability zone.
* `high_availability` - (Optional) One `high_availability` block as defined below.
* `storage` - (Optional) One `storage` block as defined below.
* `tags` - (Optional) Tags.

---

A `high_availability` block supports:

* `mode` - (Required) `SameZone` or `ZoneRedundant`.
* `standby_availability_zone` - (Optional) Standby zone.

---

A `storage` block supports:

* `size_gb` - (Optional) Storage size in GB.
* `iops` - (Optional) Storage IOPS.
* `auto_grow_enabled` - (Optional) Enable auto-grow.
* `io_scaling_enabled` - (Optional) Enable IO scaling.

## Attributes Reference

* `id` - The MySQL Flexible Server ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.DBforMySQL/flexibleServers/{name}
  ```

* `fqdn` - Simulated FQDN: `<name>.mysql.database.azure.com`.

~> **Note:** The `administrator_password` is stored in the Terraform state as sensitive data.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
