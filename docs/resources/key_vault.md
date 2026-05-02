---
subcategory: "Key Vault"
page_title: "AzureSim: azuresim_key_vault"
description: |-
  Manages a simulated Azure Key Vault.
---

# azuresim_key_vault

Manages a simulated Azure Key Vault.

This resource mimics the behavior of the [`azurerm_key_vault`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault) resource.

## Example Usage

```terraform
resource "azuresim_key_vault" "example" {
  name                       = "kv-example-001"
  resource_group_name        = azuresim_resource_group.example.name
  location                   = azuresim_resource_group.example.location
  tenant_id                  = "00000000-0000-0000-0000-000000000000"
  sku_name                   = "standard"
  enable_rbac_authorization  = true
  purge_protection_enabled   = false
  soft_delete_retention_days = 30

  tags = {
    environment = "production"
  }
}
```

## Example Usage - With Access Policy

```terraform
resource "azuresim_key_vault" "policy" {
  name                = "kv-policy-001"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  tenant_id           = azuresim_user_assigned_identity.example.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = azuresim_user_assigned_identity.example.tenant_id
    object_id = azuresim_user_assigned_identity.example.principal_id

    key_permissions    = ["Get", "List", "Create"]
    secret_permissions = ["Get", "List", "Set"]
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the Key Vault (must be globally unique in real Azure, 3-24 alphanumeric).
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `tenant_id` - (Required) Azure AD tenant ID.
* `sku_name` - (Required) `standard` or `premium`.
* `enabled_for_deployment` - (Optional) Allow VMs to retrieve secrets.
* `enabled_for_disk_encryption` - (Optional) Allow Disk Encryption to retrieve secrets and unwrap keys.
* `enabled_for_template_deployment` - (Optional) Allow Resource Manager to retrieve secrets.
* `enable_rbac_authorization` - (Optional) Use RBAC instead of access policies.
* `purge_protection_enabled` - (Optional) Enable purge protection.
* `soft_delete_retention_days` - (Optional, Computed) Soft-delete retention (7-90 days). Defaults to `90`.
* `public_network_access_enabled` - (Optional, Computed) Whether public network access is enabled. Defaults to `true`.
* `access_policy` - (Optional) Zero or more `access_policy` blocks. Mutually exclusive with `enable_rbac_authorization = true`.
* `tags` - (Optional) Tags.

---

An `access_policy` block supports:

* `tenant_id` - (Required) Azure AD tenant ID.
* `object_id` - (Required) Principal object ID.
* `application_id` - (Optional) Application ID for service principal access.
* `key_permissions` - (Optional) Permitted key operations.
* `secret_permissions` - (Optional) Permitted secret operations.
* `certificate_permissions` - (Optional) Permitted certificate operations.
* `storage_permissions` - (Optional) Permitted storage operations.

## Attributes Reference

* `id` - The Key Vault ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.KeyVault/vaults/{name}
  ```

* `vault_uri` - Simulated vault URI: `https://<name>.vault.azure.net/`
