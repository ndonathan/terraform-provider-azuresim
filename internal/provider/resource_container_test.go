package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContainerRegistry_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_container_registry" "test" {
  name                = "acr01registry"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"
  admin_enabled       = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_container_registry.test", "login_server",
						"acr01registry.azurecr.io"),
					resource.TestCheckResourceAttr("azuresim_container_registry.test", "admin_username", "acr01registry"),
					resource.TestCheckResourceAttrSet("azuresim_container_registry.test", "admin_password"),
				),
			},
		},
	})
}

func TestAccKubernetesCluster_computedFQDN(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_kubernetes_cluster" "test" {
  name                = "aks-001"
  resource_group_name = "rg"
  location            = "eastus"
  dns_prefix          = "myaks"
  kubernetes_version  = "1.30.0"
  sku_tier            = "Standard"

  default_node_pool {
    name       = "system"
    vm_size    = "Standard_D2s_v5"
    node_count = 3
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
    service_cidr   = "10.100.0.0/16"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_kubernetes_cluster.test", "node_resource_group",
						"MC_rg_aks-001_eastus"),
					resource.TestMatchResourceAttr("azuresim_kubernetes_cluster.test", "fqdn",
						regexp.MustCompile(`^myaks-deadbeef\.hcp\.eastus\.azmk8s\.io$`)),
				),
			},
		},
	})
}

func TestAccContainerAppEnvironmentAndApp(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_container_app_environment" "env" {
  name                = "cae-test"
  resource_group_name = "rg"
  location            = "eastus"
}

resource "azuresim_container_app" "app" {
  name                         = "app-svc"
  resource_group_name          = "rg"
  container_app_environment_id = azuresim_container_app_environment.env.id
  revision_mode                = "Single"

  template {
    min_replicas = 1
    max_replicas = 3
    container {
      name   = "main"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  ingress {
    external_enabled = true
    target_port      = 80
    transport        = "auto"
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_container_app_environment.env", "default_domain",
						regexp.MustCompile(`^cae-test\.eastus\.azurecontainerapps\.io$`)),
					resource.TestCheckResourceAttr("azuresim_container_app_environment.env", "static_ip_address", "203.0.113.50"),
					resource.TestCheckResourceAttr("azuresim_container_app.app", "latest_revision_fqdn",
						"app-svc.azurecontainerapps.io"),
				),
			},
		},
	})
}
