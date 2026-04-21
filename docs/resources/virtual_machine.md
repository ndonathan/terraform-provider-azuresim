---
subcategory: "Compute"
page_title: "AzureSim: azuresim_virtual_machine"
description: |-
  Manages a simulated Azure Virtual Machine.
---

# azuresim_virtual_machine

Manages a simulated Azure Virtual Machine.

This resource mimics the behavior of the [`azurerm_linux_virtual_machine`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_virtual_machine) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_virtual_machine" "example" {
  name                = "vm-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  vm_size             = "Standard_DS1_v2"
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 30
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  tags = {
    environment = "production"
  }
}
```

## Example Usage - With Network Interface

```terraform
resource "azuresim_virtual_machine" "with_nic" {
  name                = "vm-with-nic"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  vm_size             = "Standard_DS2_v2"
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"

  network_interface_ids = [
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-example/providers/Microsoft.Network/networkInterfaces/nic-example",
  ]

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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Machine. Changing this forces a new Virtual Machine to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Virtual Machine should exist. Changing this forces a new Virtual Machine to be created.

* `location` - (Required) The Azure Region where the Virtual Machine should exist. Changing this forces a new Virtual Machine to be created.

* `vm_size` - (Required) The SKU which should be used for this Virtual Machine, such as `Standard_DS1_v2`, `Standard_B2s`, or `Standard_F2`.

* `admin_username` - (Optional) The username of the local administrator used for the Virtual Machine.

* `admin_password` - (Optional) The password associated with the local administrator account.

~> **Note:** The `admin_password` is stored in the Terraform state as sensitive data.

* `network_interface_ids` - (Optional) A list of Network Interface IDs which should be attached to this Virtual Machine. These can be IDs from other simulated or real resources.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Machine.

---

An `os_disk` block supports the following:

* `caching` - (Required) The Type of Caching which should be used for the Internal OS Disk. Possible values are `None`, `ReadOnly`, and `ReadWrite`.

* `storage_account_type` - (Required) The Type of Storage Account which should back this the Internal OS Disk. Possible values are `Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS`, `StandardSSD_ZRS`, and `Premium_ZRS`.

* `disk_size_gb` - (Optional) The Size of the Internal OS Disk in GB, if you wish to vary from the size used in the image this Virtual Machine is sourced from.

---

A `source_image_reference` block supports the following:

* `publisher` - (Required) Specifies the publisher of the image used to create the Virtual Machine. For example, `Canonical` or `MicrosoftWindowsServer`.

* `offer` - (Required) Specifies the offer of the image used to create the Virtual Machine. For example, `0001-com-ubuntu-server-jammy` or `WindowsServer`.

* `sku` - (Required) Specifies the SKU of the image used to create the Virtual Machine. For example, `22_04-lts` or `2022-Datacenter`.

* `version` - (Required) Specifies the version of the image used to create the Virtual Machine. Use `latest` to use the latest available version.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine. This follows the Azure Resource Manager ID format:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Compute/virtualMachines/{name}
  ```

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
