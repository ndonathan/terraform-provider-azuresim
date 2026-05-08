---
subcategory: "Database"
page_title: "AzureSim: azuresim_cosmosdb_account"
description: |-
  Manages a simulated Azure Cosmos DB Account.
---

# azuresim_cosmosdb_account

Manages a simulated Azure Cosmos DB Account.

This resource mimics the behavior of the [`azurerm_cosmosdb_account`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cosmosdb_account) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_cosmosdb_account" "example" {
  name                = "cosmos-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location          = "eastus"
    failover_priority = 0
  }

  geo_location {
    location          = "westus2"
    failover_priority = 1
  }

  capabilities {
    name = "EnableServerless"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Account name. Must be globally unique, 3-44 lowercase alphanumeric or dashes.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `offer_type` - (Required) Always `Standard`.
* `kind` - (Optional) `GlobalDocumentDB` (default), `MongoDB`, or `Parse`.
* `automatic_failover_enabled` - (Optional) Enable automatic failover.
* `multiple_write_locations_enabled` - (Optional) Enable multi-master writes.
* `ip_range_filter` - (Optional) Comma-separated list of IP ranges allowed.
* `public_network_access_enabled` - (Optional) Allow public network access.
* `local_authentication_disabled` - (Optional) Disable key-based authentication.
* `consistency_policy` - (Required) One `consistency_policy` block as defined below.
* `geo_location` - (Required) One or more `geo_location` blocks as defined below.
* `capabilities` - (Optional) Zero or more `capabilities` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `consistency_policy` block supports:

* `consistency_level` - (Required) `Strong`, `BoundedStaleness`, `Session`, `ConsistentPrefix`, or `Eventual`.
* `max_interval_in_seconds` - (Optional) Bounded-staleness window in seconds.
* `max_staleness_prefix` - (Optional) Bounded-staleness lag in operations.

---

A `geo_location` block supports:

* `location` - (Required) Azure region.
* `failover_priority` - (Required) Failover priority (`0` = primary).
* `zone_redundant` - (Optional) Use zone redundancy.

---

A `capabilities` block supports:

* `name` - (Required) Capability name (e.g. `EnableServerless`, `EnableMongo`).

## Attributes Reference

* `id` - The Cosmos DB Account ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.DocumentDB/databaseAccounts/{name}
  ```

* `endpoint` - Simulated endpoint (`https://<name>.documents.azure.com:443/`).
* `primary_key` - Simulated primary key (sensitive).
* `secondary_key` - Simulated secondary key (sensitive).
* `primary_readonly_key` - Simulated primary readonly key (sensitive).
* `secondary_readonly_key` - Simulated secondary readonly key (sensitive).

~> **Note:** Keys are stored in the Terraform state as sensitive data. Values are static placeholders deterministically derived from RG + name and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
