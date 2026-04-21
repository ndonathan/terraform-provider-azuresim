terraform {
  required_providers {
    azuresim = {
      source = "ndonathan/azuresim"
    }
  }
}

provider "azuresim" {
  subscription_id = "12345678-1234-1234-1234-123456789012"
  tenant_id       = "87654321-4321-4321-4321-210987654321"
}

# --- Resource Group ---
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"

  tags = {
    environment = "dev"
    project     = "demo"
  }
}

# --- Virtual Network ---
resource "azuresim_virtual_network" "example" {
  name                = "vnet-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  tags = {
    environment = "dev"
  }
}

# --- Subnet ---
resource "azuresim_subnet" "example" {
  name                 = "snet-example"
  resource_group_name  = azuresim_resource_group.example.name
  virtual_network_name = azuresim_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

# --- Storage Account ---
resource "azuresim_storage_account" "example" {
  name                     = "stexample0001"
  resource_group_name      = azuresim_resource_group.example.name
  location                 = azuresim_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"

  tags = {
    environment = "dev"
  }
}

# --- Virtual Machine ---
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
    environment = "dev"
  }
}

# --- Outputs ---
output "resource_group_id" {
  value = azuresim_resource_group.example.id
}

output "vnet_id" {
  value = azuresim_virtual_network.example.id
}

output "subnet_id" {
  value = azuresim_subnet.example.id
}

output "storage_account_id" {
  value = azuresim_storage_account.example.id
}

output "storage_blob_endpoint" {
  value = azuresim_storage_account.example.primary_blob_endpoint
}

output "vm_id" {
  value = azuresim_virtual_machine.example.id
}
