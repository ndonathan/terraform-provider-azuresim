---
subcategory: "Database"
page_title: "AzureSim: azuresim_redis_cache"
description: |-
  Manages a simulated Azure Cache for Redis instance.
---

# azuresim_redis_cache

Manages a simulated Azure Cache for Redis instance.

This resource mimics the behavior of the [`azurerm_redis_cache`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/redis_cache) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_redis_cache" "example" {
  name                = "redis-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
  minimum_tls_version = "1.2"
  redis_version       = "6"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Cache name. Globally unique in real Azure.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `capacity` - (Required) Cache size (`Basic`/`Standard`: 0-6, `Premium`: 1-5).
* `family` - (Required) `C` (Basic/Standard) or `P` (Premium).
* `sku_name` - (Required) `Basic`, `Standard`, or `Premium`.
* `non_ssl_port_enabled` - (Optional) Enable the non-SSL port (6379).
* `minimum_tls_version` - (Optional) `1.0`, `1.1`, or `1.2`.
* `redis_version` - (Optional) Redis version (`4` or `6`).
* `subnet_id` - (Optional) Subnet ID (`Premium` SKU only).
* `private_static_ip_address` - (Optional) Static IP address inside the subnet.
* `zones` - (Optional) Availability zones.
* `public_network_access_enabled` - (Optional) Allow public network access.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Redis Cache ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Cache/Redis/{name}
  ```

* `hostname` - Simulated hostname (`<name>.redis.cache.windows.net`).
* `port` - Non-SSL port (`6379`).
* `ssl_port` - SSL port (`6380`).
* `primary_access_key` - Simulated primary access key (sensitive).
* `secondary_access_key` - Simulated secondary access key (sensitive).
* `primary_connection_string` - Simulated primary connection string of the form `<name>.redis.cache.windows.net:6380,password=<key>,ssl=True,abortConnect=False` (sensitive).
* `secondary_connection_string` - Simulated secondary connection string (sensitive).

~> **Note:** Keys and connection strings are stored in the Terraform state as sensitive data. Values are static placeholders deterministically derived from RG + name and should not be used for any real authentication.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
