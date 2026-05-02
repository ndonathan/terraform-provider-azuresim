package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLogAnalyticsWorkspace_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_log_analytics_workspace" "test" {
  name                = "law-test"
  resource_group_name = "rg"
  location            = "eastus"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_log_analytics_workspace.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.OperationalInsights/workspaces/law-test")),
					resource.TestCheckResourceAttr("azuresim_log_analytics_workspace.test", "sku", "PerGB2018"),
					resource.TestCheckResourceAttr("azuresim_log_analytics_workspace.test", "retention_in_days", "30"),
					resource.TestMatchResourceAttr("azuresim_log_analytics_workspace.test", "workspace_id", uuidV4Regex),
					resource.TestCheckResourceAttrSet("azuresim_log_analytics_workspace.test", "primary_shared_key"),
				),
			},
		},
	})
}

func TestAccApplicationInsights_connectionString(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_application_insights" "test" {
  name                = "appi-test"
  resource_group_name = "rg"
  location            = "eastus"
  application_type    = "web"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_application_insights.test", "instrumentation_key", uuidV4Regex),
					resource.TestMatchResourceAttr("azuresim_application_insights.test", "connection_string",
						regexp.MustCompile(`^InstrumentationKey=[0-9a-f-]+;IngestionEndpoint=https://eastus\.in\.applicationinsights\.azure\.com/`)),
				),
			},
		},
	})
}

func TestAccMonitorDiagnosticSetting_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_monitor_diagnostic_setting" "test" {
  name                       = "send-to-law"
  target_resource_id         = "/subscriptions/abc/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/kv"
  log_analytics_workspace_id = "/subscriptions/abc/resourceGroups/rg/providers/Microsoft.OperationalInsights/workspaces/law"

  enabled_log {
    category_group = "allLogs"
  }

  metric {
    category = "AllMetrics"
    enabled  = true
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_monitor_diagnostic_setting.test", "id",
					"/subscriptions/abc/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/kv|send-to-law"),
			},
		},
	})
}

func TestAccMonitorActionGroup_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_monitor_action_group" "test" {
  name                = "ag-oncall"
  resource_group_name = "rg"
  short_name          = "oncall"

  email_receiver {
    name          = "platform"
    email_address = "platform@example.com"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_monitor_action_group.test", "id",
					arm("/resourceGroups/rg/providers/microsoft.insights/actionGroups/ag-oncall")),
			},
		},
	})
}

func TestAccMonitorMetricAlert_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_monitor_metric_alert" "test" {
  name                = "ma-cpu"
  resource_group_name = "rg"
  scopes              = ["/subscriptions/abc/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm"]
  severity            = 2
  frequency           = "PT1M"
  window_size         = "PT5M"

  criteria {
    metric_namespace = "Microsoft.Compute/virtualMachines"
    metric_name      = "Percentage CPU"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 80
  }

  action {
    action_group_id = "/subscriptions/abc/resourceGroups/rg/providers/microsoft.insights/actionGroups/ag-oncall"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_monitor_metric_alert.test", "id",
					arm("/resourceGroups/rg/providers/microsoft.insights/metricAlerts/ma-cpu")),
			},
		},
	})
}

func TestAccRecoveryServicesVault_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_recovery_services_vault" "test" {
  name                = "rsv-001"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"
  storage_mode_type   = "GeoRedundant"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_recovery_services_vault.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.RecoveryServices/vaults/rsv-001")),
			},
		},
	})
}

func TestAccApplicationGateway_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_application_gateway" "test" {
  name                = "agw-001"
  resource_group_name = "rg"
  location            = "eastus"

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "appgw-ipconf"
    subnet_id = "/subnets/agw"
  }

  frontend_port {
    name = "http"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "/pip"
  }

  backend_address_pool {
    name = "pool1"
  }

  backend_http_settings {
    name = "settings1"
  }

  http_listener {
    name = "listener1"
  }

  request_routing_rule {
    name = "rule1"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_application_gateway.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/applicationGateways/agw-001")),
			},
		},
	})
}

func TestAccFirewall_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_firewall" "test" {
  name                = "afw-001"
  resource_group_name = "rg"
  location            = "eastus"
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = "/AzureFirewallSubnet"
    public_ip_address_id = "/pip"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_firewall.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/azureFirewalls/afw-001")),
			},
		},
	})
}

func TestAccAPIManagement_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_api_management" "test" {
  name                = "apim-001"
  resource_group_name = "rg"
  location            = "eastus"
  publisher_name      = "Acme Corp"
  publisher_email     = "ops@example.com"
  sku_name            = "Standard_2"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_api_management.test", "gateway_url",
						"https://apim-001.azure-api.net"),
					resource.TestCheckResourceAttr("azuresim_api_management.test", "developer_portal_url",
						"https://apim-001.developer.azure-api.net"),
				),
			},
		},
	})
}

func TestAccFrontDoorProfile_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_cdn_frontdoor_profile" "test" {
  name                = "afd-001"
  resource_group_name = "rg"
  sku_name            = "Standard_AzureFrontDoor"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_cdn_frontdoor_profile.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Cdn/profiles/afd-001")),
			},
		},
	})
}

func TestAccPrivateEndpoint_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_private_endpoint" "test" {
  name                = "pe-001"
  resource_group_name = "rg"
  location            = "eastus"
  subnet_id           = "/subnets/pe"

  private_service_connection {
    name                           = "psc"
    private_connection_resource_id = "/sa"
    is_manual_connection           = false
    subresource_names              = ["blob"]
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_private_endpoint.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Network/privateEndpoints/pe-001")),
			},
		},
	})
}
