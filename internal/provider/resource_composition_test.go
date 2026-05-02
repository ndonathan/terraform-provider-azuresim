package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccComposition_kitchenSink stitches a realistic Azure footprint together:
// RG → VNet → Subnet → NSG (with rule + association) → Public IP → NIC → VM,
// plus a Storage Account / Container / Blob, a Key Vault with a Secret, a
// User-Assigned Identity, and a Role Assignment. The goal is to catch
// regressions in cross-resource interpolation: every dependent resource must
// be able to consume the parent's computed `id` (or other computed attribute)
// without producing inconsistent-result errors.
func TestAccComposition_kitchenSink(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-kitchen-sink"
  location = "eastus"
  tags     = { env = "test" }
}

# --- Networking ---

resource "azuresim_virtual_network" "vnet" {
  name                = "vnet"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "azuresim_subnet" "snet" {
  name                 = "snet-app"
  resource_group_name  = azuresim_resource_group.rg.name
  virtual_network_name = azuresim_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azuresim_network_security_group" "nsg" {
  name                = "nsg"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location

  security_rule {
    name                       = "AllowSSH"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azuresim_subnet_network_security_group_association" "snet_nsg" {
  subnet_id                 = azuresim_subnet.snet.id
  network_security_group_id = azuresim_network_security_group.nsg.id
}

resource "azuresim_public_ip" "pip" {
  name                = "pip"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "kitchen-sink"
}

resource "azuresim_network_interface" "nic" {
  name                = "nic"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azuresim_subnet.snet.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.1.10"
    public_ip_address_id          = azuresim_public_ip.pip.id
  }
}

# --- Compute ---

resource "azuresim_virtual_machine" "vm" {
  name                  = "vm"
  resource_group_name   = azuresim_resource_group.rg.name
  location              = azuresim_resource_group.rg.location
  vm_size               = "Standard_DS1_v2"
  admin_username        = "azureuser"
  admin_password        = "` + testFixturePassword + `"
  network_interface_ids = [azuresim_network_interface.nic.id]

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
}

# --- Identity / RBAC ---

resource "azuresim_user_assigned_identity" "uai" {
  name                = "uai-app"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
}

resource "azuresim_role_assignment" "rbac" {
  scope                = azuresim_resource_group.rg.id
  role_definition_name = "Reader"
  principal_id         = azuresim_user_assigned_identity.uai.principal_id
  principal_type       = "ServicePrincipal"
}

# --- Storage ---

resource "azuresim_storage_account" "sa" {
  name                     = "stkitchensink01"
  resource_group_name      = azuresim_resource_group.rg.name
  location                 = azuresim_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azuresim_storage_container" "container" {
  name                  = "data"
  storage_account_name  = azuresim_storage_account.sa.name
  container_access_type = "private"
}

resource "azuresim_storage_blob" "blob" {
  name                   = "config.json"
  storage_account_name   = azuresim_storage_account.sa.name
  storage_container_name = azuresim_storage_container.container.name
  type                   = "BlockBlob"
  source_content         = "{}"
}

# --- Key Vault ---

resource "azuresim_key_vault" "kv" {
  name                = "kvkitchensink01"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_key_vault_secret" "secret" {
  name         = "db-password"
  key_vault_id = azuresim_key_vault.kv.id
  value        = "` + testFixtureSecretValue + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Top-of-graph IDs are well-formed.
					resource.TestCheckResourceAttr("azuresim_resource_group.rg", "id",
						arm("/resourceGroups/rg-kitchen-sink")),

					// Subnet–NSG association uses the subnet's ID as its own ID.
					resource.TestCheckResourceAttrPair(
						"azuresim_subnet_network_security_group_association.snet_nsg", "id",
						"azuresim_subnet.snet", "id",
					),

					// VM is wired to the NIC.
					resource.TestCheckResourceAttrPair(
						"azuresim_virtual_machine.vm", "network_interface_ids.0",
						"azuresim_network_interface.nic", "id",
					),

					// Role assignment scope and principal are interpolated correctly.
					resource.TestCheckResourceAttrPair(
						"azuresim_role_assignment.rbac", "scope",
						"azuresim_resource_group.rg", "id",
					),
					resource.TestCheckResourceAttrPair(
						"azuresim_role_assignment.rbac", "principal_id",
						"azuresim_user_assigned_identity.uai", "principal_id",
					),

					// Storage chain end-to-end.
					resource.TestCheckResourceAttr("azuresim_storage_blob.blob", "id",
						"https://stkitchensink01.blob.core.windows.net/data/config.json"),

					// Key Vault Secret URI references the parent vault.
					resource.TestCheckResourceAttr("azuresim_key_vault_secret.secret", "versionless_id",
						"https://kvkitchensink01.vault.azure.net/secrets/db-password"),
				),
			},
		},
	})
}
