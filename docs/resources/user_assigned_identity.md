---
subcategory: "Identity"
page_title: "AzureSim: azuresim_user_assigned_identity"
description: |-
  Manages a simulated Azure User-Assigned Managed Identity.
---

# azuresim_user_assigned_identity

Manages a simulated Azure User-Assigned Managed Identity.

This resource mimics the behavior of the [`azurerm_user_assigned_identity`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/user_assigned_identity) resource.

## Example Usage

```terraform
resource "azuresim_user_assigned_identity" "example" {
  name                = "uai-app"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the identity.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Managed Identity ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{name}
  ```

* `principal_id` - Simulated service principal object ID. Deterministically derived from the resource group + name.
* `client_id` - Simulated client (application) ID. Deterministically derived from the resource group + name.
* `tenant_id` - Tenant ID. Currently always the zero UUID.
