---
subcategory: "Compute"
page_title: "AzureSim: azuresim_windows_virtual_machine"
description: |-
  Manages a simulated Azure Windows Virtual Machine.
---

# azuresim_windows_virtual_machine

Manages a simulated Azure Windows Virtual Machine.

This resource mimics the behavior of the [`azurerm_windows_virtual_machine`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_virtual_machine) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_network_interface" "example" {
  name                = "nic-winvm"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azuresim_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azuresim_windows_virtual_machine" "example" {
  name                = "vm-win"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  size                = "Standard_D2s_v5"
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"

  network_interface_ids = [azuresim_network_interface.example.id]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
    disk_size_gb         = 128
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition"
    version   = "latest"
  }

  license_type        = "Windows_Server"
  patch_mode          = "AutomaticByPlatform"
  hotpatching_enabled = true

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) VM name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `size` - (Required) VM size (e.g. `Standard_D2s_v5`).
* `admin_username` - (Required) Local administrator username.
* `admin_password` - (Required) Local administrator password (sensitive).
* `computer_name` - (Optional) Computer name. Defaults to the VM name when omitted.
* `network_interface_ids` - (Required) List of NIC IDs to attach. The first is treated as primary.
* `os_disk` - (Required) One `os_disk` block as defined below.
* `source_image_reference` - (Required) One `source_image_reference` block as defined below.
* `license_type` - (Optional) `None`, `Windows_Client`, or `Windows_Server`.
* `hotpatching_enabled` - (Optional) Enable hotpatching (Windows Server 2022 Datacenter Azure Edition only).
* `patch_mode` - (Optional) `Manual`, `AutomaticByOS`, or `AutomaticByPlatform`.
* `zone` - (Optional, ForceNew) Availability zone.
* `tags` - (Optional) Tags.

~> **Note:** The `admin_password` is stored in the Terraform state as sensitive data.

---

An `os_disk` block supports:

* `caching` - (Required) `None`, `ReadOnly`, or `ReadWrite`.
* `storage_account_type` - (Required) `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS`, etc.
* `disk_size_gb` - (Optional) Disk size in GB.
* `name` - (Optional) Disk name.

---

A `source_image_reference` block supports:

* `publisher` - (Required) e.g. `MicrosoftWindowsServer`.
* `offer` - (Required) e.g. `WindowsServer`.
* `sku` - (Required) e.g. `2022-datacenter-azure-edition`.
* `version` - (Required) e.g. `latest`.

## Attributes Reference

* `id` - The Virtual Machine ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Compute/virtualMachines/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
