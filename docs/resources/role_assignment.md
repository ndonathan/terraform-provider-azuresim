---
subcategory: "Authorization"
page_title: "AzureSim: azuresim_role_assignment"
description: |-
  Manages a simulated Azure RBAC Role Assignment.
---

# azuresim_role_assignment

Manages a simulated Azure RBAC Role Assignment.

This resource mimics the behavior of the [`azurerm_role_assignment`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_user_assigned_identity" "example" {
  name                = "uai-app"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
}

resource "azuresim_role_assignment" "example" {
  scope                = azuresim_resource_group.example.id
  role_definition_name = "Contributor"
  principal_id         = azuresim_user_assigned_identity.example.principal_id
  principal_type       = "ServicePrincipal"
  description          = "Grants the app identity Contributor on the resource group"
}
```

## Argument Reference

* `name` - (Optional, Computed, ForceNew) Assignment name (UUID). Generated deterministically when omitted.
* `scope` - (Required, ForceNew) Scope at which the role applies (subscription, resource group, or resource ID).
* `role_definition_id` - (Optional, ForceNew) Role definition ID. Mutually exclusive with `role_definition_name`.
* `role_definition_name` - (Optional, ForceNew) Role definition name (e.g. `Contributor`, `Reader`). Mutually exclusive with `role_definition_id`.
* `principal_id` - (Required, ForceNew) Object ID of the principal.
* `principal_type` - (Optional, Computed) Principal type (`User`, `Group`, or `ServicePrincipal`). Defaults to `ServicePrincipal` when omitted.
* `description` - (Optional) Description for the assignment.
* `condition` - (Optional) ABAC condition expression.
* `condition_version` - (Optional) Condition version (e.g. `2.0`).
* `skip_service_principal_aad_check` - (Optional) Skip AAD propagation check (no-op in this simulator).
* `delegated_managed_identity_resource_id` - (Optional) Delegated managed identity resource ID.

## Attributes Reference

* `id` - The Role Assignment ID:

  ```
  {scope}/providers/Microsoft.Authorization/roleAssignments/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
