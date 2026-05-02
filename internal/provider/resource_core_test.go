package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Builds a `/subscriptions/<sub>/...` prefix with the test subscription ID,
// so individual tests can assert exact ARM IDs without string-concatenating
// the constant in every check.
func arm(suffix string) string {
	return fmt.Sprintf("/subscriptions/%s%s", testSubscriptionID, suffix)
}

func TestAccResourceGroup_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "test" {
  name     = "rg-test"
  location = "eastus"
  tags     = { env = "dev" }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "id", arm("/resourceGroups/rg-test")),
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "name", "rg-test"),
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "location", "eastus"),
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "tags.env", "dev"),
				),
			},
			{
				// Tag-only update preserves ID.
				Config: providerConfig + `
resource "azuresim_resource_group" "test" {
  name     = "rg-test"
  location = "eastus"
  tags     = { env = "prod", owner = "platform" }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "id", arm("/resourceGroups/rg-test")),
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "tags.env", "prod"),
					resource.TestCheckResourceAttr("azuresim_resource_group.test", "tags.owner", "platform"),
				),
			},
		},
	})
}

func TestAccResourceGroup_renameForcesReplace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "test" {
  name     = "rg-original"
  location = "eastus"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_resource_group.test", "id", arm("/resourceGroups/rg-original")),
			},
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "test" {
  name     = "rg-renamed"
  location = "eastus"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_resource_group.test", "id", arm("/resourceGroups/rg-renamed")),
			},
		},
	})
}

func TestAccVirtualNetwork_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-vnet"
  location = "eastus"
}

resource "azuresim_virtual_network" "test" {
  name                = "vnet-test"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  address_space       = ["10.0.0.0/16", "172.16.0.0/12"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_virtual_network.test", "id", arm("/resourceGroups/rg-vnet/providers/Microsoft.Network/virtualNetworks/vnet-test")),
					resource.TestCheckResourceAttr("azuresim_virtual_network.test", "address_space.#", "2"),
					resource.TestCheckResourceAttr("azuresim_virtual_network.test", "address_space.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("azuresim_virtual_network.test", "dns_servers.0", "10.0.0.4"),
				),
			},
		},
	})
}

func TestAccSubnet_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_resource_group" "rg" {
  name     = "rg-subnet"
  location = "eastus"
}

resource "azuresim_virtual_network" "vnet" {
  name                = "vnet-subnet"
  resource_group_name = azuresim_resource_group.rg.name
  location            = azuresim_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "azuresim_subnet" "test" {
  name                 = "snet"
  resource_group_name  = azuresim_resource_group.rg.name
  virtual_network_name = azuresim_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_subnet.test", "id",
						arm("/resourceGroups/rg-subnet/providers/Microsoft.Network/virtualNetworks/vnet-subnet/subnets/snet")),
					resource.TestCheckResourceAttr("azuresim_subnet.test", "address_prefixes.0", "10.0.1.0/24"),
				),
			},
		},
	})
}
