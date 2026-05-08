---
subcategory: "Network"
page_title: "AzureSim: azuresim_application_gateway"
description: |-
  Manages a simulated Azure Application Gateway.
---

# azuresim_application_gateway

Manages a simulated Azure Application Gateway.

This resource mimics the behavior of the [`azurerm_application_gateway`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/application_gateway) resource. It manages all state within Terraform's state file and does not make any API calls to Azure. The schema captures the major top-level blocks; detailed block attributes are intentionally minimal in this simulator.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_application_gateway" "example" {
  name                = "agw-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  enable_http2        = true

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "gateway-ip-config"
    subnet_id = azuresim_subnet.example.id
  }

  frontend_port {
    name = "http"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "public"
    public_ip_address_id = azuresim_public_ip.example.id
  }

  backend_address_pool {
    name = "backend-pool"
  }

  backend_http_settings {
    name = "http-settings"
  }

  http_listener {
    name = "http-listener"
  }

  request_routing_rule {
    name = "rule-1"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the Application Gateway.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `zones` - (Optional) List of availability zones.
* `enable_http2` - (Optional) Whether HTTP/2 is enabled.
* `fips_enabled` - (Optional) Whether FIPS mode is enabled.
* `firewall_policy_id` - (Optional) WAF policy ID to associate.
* `sku` - (Required) One `sku` block as defined below.
* `gateway_ip_configuration` - (Required) One or more `gateway_ip_configuration` blocks as defined below.
* `frontend_port` - (Required) One or more `frontend_port` blocks as defined below.
* `frontend_ip_configuration` - (Required) One or more `frontend_ip_configuration` blocks as defined below.
* `backend_address_pool` - (Required) One or more `backend_address_pool` blocks as defined below.
* `backend_http_settings` - (Required) One or more `backend_http_settings` blocks as defined below.
* `http_listener` - (Required) One or more `http_listener` blocks as defined below.
* `request_routing_rule` - (Required) One or more `request_routing_rule` blocks as defined below.
* `tags` - (Optional) Tags.

---

A `sku` block supports:

* `name` - (Required) SKU name (e.g. `Standard_v2`, `WAF_v2`).
* `tier` - (Required) SKU tier.
* `capacity` - (Optional) Fixed capacity (instance count).

---

A `gateway_ip_configuration` block supports:

* `name` - (Required) Configuration name.
* `subnet_id` - (Required) Subnet ID where the gateway is deployed.

---

A `frontend_port` block supports:

* `name` - (Required) Port name.
* `port` - (Required) Port number.

---

A `frontend_ip_configuration` block supports:

* `name` - (Required) Configuration name.
* `public_ip_address_id` - (Optional) Public IP ID.
* `subnet_id` - (Optional) Subnet ID for a private listener.
* `private_ip_address` - (Optional) Static private IP.
* `private_ip_address_allocation` - (Optional) `Static` or `Dynamic`.

---

A `backend_address_pool`, `backend_http_settings`, `http_listener`, or `request_routing_rule` block supports:

* `name` - (Required) Block name. Detailed sub-attributes are not modeled in this simulator.

## Attributes Reference

* `id` - The Application Gateway ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Network/applicationGateways/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
