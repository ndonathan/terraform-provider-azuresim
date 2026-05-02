---
subcategory: "Compute"
page_title: "AzureSim: azuresim_managed_disk"
description: |-
  Manages a simulated Azure Managed Disk.
---

# azuresim_managed_disk

Manages a simulated Azure Managed Disk.

This resource mimics the behavior of the [`azurerm_managed_disk`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/managed_disk) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_managed_disk" "example" {
  name                 = "disk-data-01"
  resource_group_name  = azuresim_resource_group.example.name
  location             = azuresim_resource_group.example.location
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 256

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Name of the Managed Disk.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `storage_account_type` - (Required) `Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS`, `UltraSSD_LRS`, `Premium_ZRS`, or `StandardSSD_ZRS`.
* `create_option` - (Required, ForceNew) `Empty`, `Copy`, `FromImage`, `Import`, `Restore`, or `Upload`.
* `disk_size_gb` - (Optional, Computed) Disk size in GB. Defaults to `30`.
* `source_uri` - (Optional) Source blob URI when `create_option = "Import"`.
* `source_resource_id` - (Optional) Source resource ID when `create_option = "Copy"` or `"Restore"`.
* `image_reference_id` - (Optional) Image reference ID when `create_option = "FromImage"`.
* `os_type` - (Optional) `Linux` or `Windows`.
* `tier` - (Optional) Performance tier (e.g. `P30`).
* `max_shares` - (Optional) Number of VMs that can attach the disk simultaneously.
* `zone` - (Optional, ForceNew) Availability zone.
* `network_access_policy` - (Optional) `AllowAll`, `AllowPrivate`, or `DenyAll`.
* `public_network_access_enabled` - (Optional) Whether public network access is allowed.
* `disk_iops_read_write` - (Optional) Provisioned IOPS (UltraSSD/PremiumV2 only).
* `disk_mbps_read_write` - (Optional) Provisioned bandwidth in MB/s (UltraSSD/PremiumV2 only).
* `tags` - (Optional) Tags.

## Attributes Reference

* `id` - The Managed Disk ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Compute/disks/{name}
  ```
