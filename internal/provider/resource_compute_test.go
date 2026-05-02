package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-vm"
  location = "eastus"
}

resource "azuresim_virtual_machine" "test" {
  name                = "vm-test"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  vm_size             = "Standard_DS1_v2"
  admin_username      = "azureuser"
  admin_password      = "` + testFixturePassword + `"

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
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_virtual_machine.test", "id",
						arm("/resourceGroups/rg-vm/providers/Microsoft.Compute/virtualMachines/vm-test")),
					resource.TestCheckResourceAttr("azuresim_virtual_machine.test", "os_disk.0.caching", "ReadWrite"),
					resource.TestCheckResourceAttr("azuresim_virtual_machine.test", "source_image_reference.0.publisher", "Canonical"),
				),
			},
		},
	})
}

func TestAccWindowsVirtualMachine_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-winvm"
  location = "eastus"
}

resource "azuresim_windows_virtual_machine" "test" {
  name                  = "vm-win"
  resource_group_name   = azuresim_resource_group.rg.name
  location              = azuresim_resource_group.rg.location
  size                  = "Standard_D2s_v5"
  admin_username        = "winadmin"
  admin_password        = "` + testFixturePassword + `"
  network_interface_ids = []

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "StandardSSD_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition"
    version   = "latest"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_windows_virtual_machine.test", "id",
					arm("/resourceGroups/rg-winvm/providers/Microsoft.Compute/virtualMachines/vm-win")),
			},
		},
	})
}

func TestAccManagedDisk_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-disk"
  location = "eastus"
}

resource "azuresim_managed_disk" "test" {
  name                 = "disk-data"
  resource_group_name  = azuresim_resource_group.rg.name
  location             = azuresim_resource_group.rg.location
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 256
  os_type              = "Linux"
  zone                 = "1"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_managed_disk.test", "id",
						arm("/resourceGroups/rg-disk/providers/Microsoft.Compute/disks/disk-data")),
					resource.TestCheckResourceAttr("azuresim_managed_disk.test", "disk_size_gb", "256"),
				),
			},
			{
				// Disk size update should not force replacement.
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-disk"
  location = "eastus"
}

resource "azuresim_managed_disk" "test" {
  name                 = "disk-data"
  resource_group_name  = azuresim_resource_group.rg.name
  location             = azuresim_resource_group.rg.location
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 512
  os_type              = "Linux"
  zone                 = "1"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_managed_disk.test", "disk_size_gb", "512"),
			},
		},
	})
}

func TestAccManagedDisk_defaultDiskSize(t *testing.T) {
	// When disk_size_gb is omitted, the resource should compute a default of 30.
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_managed_disk" "test" {
  name                 = "disk-default"
  resource_group_name  = "rg"
  location             = "eastus"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_managed_disk.test", "disk_size_gb", "30"),
			},
		},
	})
}

func TestAccLinuxVMSS_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_linux_virtual_machine_scale_set" "test" {
  name                = "vmss-linux"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard_D2s_v5"
  instances           = 3
  admin_username      = "azureuser"
  admin_password      = "` + testFixturePassword + `"

  network_interface {
    name    = "nic"
    primary = true
    ip_configuration {
      name      = "internal"
      subnet_id = "/subnets/x"
      primary   = true
    }
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "StandardSSD_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}`,
				Check: resource.TestMatchResourceAttr("azuresim_linux_virtual_machine_scale_set.test", "id",
					regexp.MustCompile(`/Microsoft\.Compute/virtualMachineScaleSets/vmss-linux$`)),
			},
		},
	})
}
