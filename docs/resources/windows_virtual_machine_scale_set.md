---
subcategory: "Compute"
page_title: "AzureSim: azuresim_windows_virtual_machine_scale_set"
description: |-
  Manages a simulated Azure Windows Virtual Machine Scale Set.
---

# azuresim_windows_virtual_machine_scale_set

Manages a simulated Azure Windows Virtual Machine Scale Set.

This resource mimics the behavior of the [`azurerm_windows_virtual_machine_scale_set`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_virtual_machine_scale_set) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_windows_virtual_machine_scale_set" "example" {
  name                = "vmss-win"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  sku                 = "Standard_D2s_v5"
  instances           = 3
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"
  upgrade_mode        = "Manual"

  network_interface {
    name    = "primary"
    primary = true

    ip_configuration {
      name      = "internal"
      subnet_id = azuresim_subnet.example.id
      primary   = true
    }
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) VMSS name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `sku` - (Required) VM SKU.
* `instances` - (Required) Initial instance count.
* `admin_username` - (Required) Admin username.
* `admin_password` - (Optional) Admin password (sensitive).
* `network_interface` - (Required) One or more `network_interface` blocks as defined below.
* `os_disk` - (Required) One `os_disk` block as defined below.
* `source_image_reference` - (Required) One `source_image_reference` block as defined below.
* `zones` - (Optional) Availability zones.
* `upgrade_mode` - (Optional) `Manual`, `Automatic`, or `Rolling`.
* `overprovision` - (Optional) Whether to overprovision VMs during scaling.
* `tags` - (Optional) Tags.

---

A `network_interface` block supports:

* `name` - (Required) NIC profile name.
* `primary` - (Optional) Whether this is the primary NIC.
* `ip_configuration` - (Required) One or more `ip_configuration` blocks.

An `ip_configuration` block (under `network_interface`) supports:

* `name` - (Required) Configuration name.
* `subnet_id` - (Required) Subnet ID.
* `primary` - (Optional) Primary IP config.
* `load_balancer_backend_address_pool_ids` - (Optional) Load balancer backend pool IDs.

---

An `os_disk` block supports:

* `caching` - (Required) `None`, `ReadOnly`, or `ReadWrite`.
* `storage_account_type` - (Required) Storage type.
* `disk_size_gb` - (Optional) Disk size in GB.

---

A `source_image_reference` block supports:

* `publisher` - (Required) Image publisher.
* `offer` - (Required) Image offer.
* `sku` - (Required) Image SKU.
* `version` - (Required) Image version.

## Attributes Reference

* `id` - The VMSS ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Compute/virtualMachineScaleSets/{name}
  ```

~> **Note:** The `admin_password` is stored in the Terraform state as sensitive data.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
