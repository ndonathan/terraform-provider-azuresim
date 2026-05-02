package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkSecurityGroup_withRules(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_network_security_group" "test" {
  name                = "nsg-test"
  resource_group_name = "rg"
  location            = "eastus"

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

  security_rule {
    name                         = "AllowWeb"
    priority                     = 110
    direction                    = "Inbound"
    access                       = "Allow"
    protocol                     = "Tcp"
    source_port_range            = "*"
    destination_port_ranges      = ["80", "443"]
    source_address_prefixes      = ["10.0.0.0/8"]
    destination_address_prefix   = "*"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_network_security_group.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/nsg-test")),
					resource.TestCheckResourceAttr("azuresim_network_security_group.test", "security_rule.#", "2"),
					resource.TestCheckResourceAttr("azuresim_network_security_group.test", "security_rule.0.priority", "100"),
					resource.TestCheckResourceAttr("azuresim_network_security_group.test", "security_rule.1.destination_port_ranges.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkSecurityRule_standalone(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_network_security_rule" "test" {
  name                        = "AllowDB"
  resource_group_name         = "rg"
  network_security_group_name = "nsg"
  protocol                    = "Tcp"
  access                      = "Allow"
  priority                    = 200
  direction                   = "Inbound"
  source_port_range           = "*"
  destination_port_range      = "5432"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_network_security_rule.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/nsg/securityRules/AllowDB")),
			},
		},
	})
}

func TestAccPublicIP_computedAddressAndFQDN(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_public_ip" "test" {
  name                = "pip-test"
  resource_group_name = "rg"
  location            = "eastus"
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "myapp"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/pip-test")),
					resource.TestMatchResourceAttr("azuresim_public_ip.test", "ip_address",
						regexp.MustCompile(`^203\.0\.113\.\d+$`)),
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "fqdn", "myapp.eastus.cloudapp.azure.com"),
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "sku", "Standard"),
				),
			},
		},
	})
}

func TestAccPublicIP_defaultSKU(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_public_ip" "test" {
  name                = "pip-default"
  resource_group_name = "rg"
  location            = "eastus"
  allocation_method   = "Dynamic"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "sku", "Basic"),
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "sku_tier", "Regional"),
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "ip_version", "IPv4"),
					resource.TestCheckResourceAttr("azuresim_public_ip.test", "fqdn", ""),
				),
			},
		},
	})
}

func TestAccPublicIP_IPv6(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_public_ip" "test" {
  name                = "pip-v6"
  resource_group_name = "rg"
  location            = "eastus"
  allocation_method   = "Static"
  ip_version          = "IPv6"
}`,
				Check: resource.TestMatchResourceAttr("azuresim_public_ip.test", "ip_address",
					regexp.MustCompile(`^2001:db8::`)),
			},
		},
	})
}

func TestAccNetworkInterface_basicAndComputed(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_network_interface" "test" {
  name                = "nic-test"
  resource_group_name = "rg"
  location            = "eastus"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "/subnets/x"
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.1.42"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_network_interface.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic-test")),
					resource.TestMatchResourceAttr("azuresim_network_interface.test", "mac_address",
						regexp.MustCompile(`^00-15-5D-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}$`)),
					resource.TestCheckResourceAttr("azuresim_network_interface.test", "private_ip_address", "10.0.1.42"),
					resource.TestCheckResourceAttr("azuresim_network_interface.test", "private_ip_addresses.#", "1"),
					resource.TestCheckResourceAttr("azuresim_network_interface.test", "internal_domain_name_suffix", "internal.cloudapp.net"),
					resource.TestCheckResourceAttr("azuresim_network_interface.test", "ip_configuration.0.primary", "true"),
				),
			},
		},
	})
}

func TestAccNetworkInterface_dynamicGetsSimulatedIP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_network_interface" "test" {
  name                = "nic-dyn"
  resource_group_name = "rg"
  location            = "eastus"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "/subnets/x"
    private_ip_address_allocation = "Dynamic"
  }
}`,
				Check: resource.TestMatchResourceAttr("azuresim_network_interface.test", "private_ip_address",
					regexp.MustCompile(`^10\.0\.0\.\d+$`)),
			},
		},
	})
}

func TestAccSubnetNSGAssociation_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_subnet_network_security_group_association" "test" {
  subnet_id                 = "/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/v/subnets/s"
  network_security_group_id = "/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/n"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_subnet_network_security_group_association.test", "id",
					"/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/v/subnets/s"),
			},
		},
	})
}

func TestAccRouteTable_withRoutes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_route_table" "test" {
  name                = "rt-test"
  resource_group_name = "rg"
  location            = "eastus"

  route {
    name           = "to-firewall"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VirtualAppliance"
    next_hop_in_ip_address = "10.0.100.4"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_route_table.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Network/routeTables/rt-test")),
					resource.TestCheckResourceAttr("azuresim_route_table.test", "route.#", "1"),
				),
			},
		},
	})
}

func TestAccPrivateDNSZone_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_private_dns_zone" "test" {
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = "rg"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_private_dns_zone.test", "max_number_of_record_sets", "25000"),
					resource.TestCheckResourceAttr("azuresim_private_dns_zone.test", "max_number_of_virtual_network_links", "1000"),
				),
			},
		},
	})
}

func TestAccVirtualNetworkPeering_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_virtual_network_peering" "test" {
  name                      = "peer-a-to-b"
  resource_group_name       = "rg"
  virtual_network_name      = "vnet-a"
  remote_virtual_network_id = "/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet-b"
  allow_virtual_network_access = true
}`,
				Check: resource.TestCheckResourceAttr("azuresim_virtual_network_peering.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet-a/virtualNetworkPeerings/peer-a-to-b")),
			},
		},
	})
}

func TestAccLoadBalancer_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_lb" "test" {
  name                = "lb-test"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "/pip"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_lb.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/loadBalancers/lb-test")),
			},
		},
	})
}

func TestAccNATGateway_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_nat_gateway" "test" {
  name                    = "nat-test"
  resource_group_name     = "rg"
  location                = "eastus"
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_nat_gateway.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Network/natGateways/nat-test")),
					resource.TestMatchResourceAttr("azuresim_nat_gateway.test", "resource_guid", uuidV4Regex),
				),
			},
		},
	})
}

func TestAccBastionHost_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_bastion_host" "test" {
  name                = "bast-test"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"

  ip_configuration {
    name                 = "config"
    subnet_id            = "/AzureBastionSubnet"
    public_ip_address_id = "/pip"
  }
}`,
				Check: resource.TestMatchResourceAttr("azuresim_bastion_host.test", "dns_name",
					regexp.MustCompile(`^bst-[0-9a-f]{8}\.bastion\.azure\.com$`)),
			},
		},
	})
}
