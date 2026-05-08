---
subcategory: "Integration"
page_title: "AzureSim: azuresim_api_management"
description: |-
  Manages a simulated Azure API Management instance.
---

# azuresim_api_management

Manages a simulated Azure API Management instance.

This resource mimics the behavior of the [`azurerm_api_management`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/api_management) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_api_management" "example" {
  name                = "apim-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  publisher_name      = "Example Org"
  publisher_email     = "apiteam@example.com"
  sku_name            = "Developer_1"

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the API Management instance. In a real Azure environment this must be globally unique.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `publisher_name` - (Required) Publisher (organization) name shown in the developer portal.
* `publisher_email` - (Required) Publisher email used for notifications.
* `sku_name` - (Required) SKU (e.g. `Developer_1`, `Standard_2`, `Premium_4`, `Consumption_0`, `BasicV2_1`, `StandardV2_1`).
* `zones` - (Optional) List of availability zones for the deployment.
* `virtual_network_type` - (Optional) `None`, `External`, or `Internal`.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The API Management ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ApiManagement/service/{name}
  ```

* `gateway_url` - Simulated gateway URL: `https://<name>.azure-api.net`
* `portal_url` - Simulated legacy publisher portal URL: `https://<name>.portal.azure-api.net`
* `management_api_url` - Simulated management API URL: `https://<name>.management.azure-api.net`
* `developer_portal_url` - Simulated developer portal URL: `https://<name>.developer.azure-api.net`
* `public_ip_addresses` - Simulated list of public IPs assigned to the gateway.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
