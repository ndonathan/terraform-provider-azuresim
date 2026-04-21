---
subcategory: "Base"
page_title: "AzureSim: azuresim_resource_group"
description: |-
  Manages a simulated Azure Resource Group.
---

# azuresim_resource_group

Manages a simulated Azure Resource Group.

This resource mimics the behavior of the [`azurerm_resource_group`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"

  tags = {
    environment = "production"
    department  = "engineering"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the Resource Group. Changing this forces a new Resource Group to be created.

* `location` - (Required) The Azure Region where the Resource Group should exist. Changing this forces a new Resource Group to be created. A full list of Azure Regions can be referenced when targeting a real environment (e.g. `eastus`, `westus2`, `westeurope`).

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
