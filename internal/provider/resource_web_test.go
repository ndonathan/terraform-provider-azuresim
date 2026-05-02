package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServicePlan_LinuxIsReserved(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_service_plan" "test" {
  name                = "asp-linux"
  resource_group_name = "rg"
  location            = "eastus"
  os_type             = "Linux"
  sku_name            = "P1v3"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_service_plan.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Web/serverfarms/asp-linux")),
					resource.TestCheckResourceAttr("azuresim_service_plan.test", "reserved", "true"),
					resource.TestCheckResourceAttr("azuresim_service_plan.test", "kind", "linux"),
				),
			},
		},
	})
}

func TestAccServicePlan_WindowsContainer(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_service_plan" "test" {
  name                = "asp-winc"
  resource_group_name = "rg"
  location            = "eastus"
  os_type             = "WindowsContainer"
  sku_name            = "P1v3"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_service_plan.test", "reserved", "false"),
					resource.TestCheckResourceAttr("azuresim_service_plan.test", "kind", "windows,container"),
				),
			},
		},
	})
}

func TestAccLinuxWebApp_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_service_plan" "plan" {
  name                = "asp-app"
  resource_group_name = "rg"
  location            = "eastus"
  os_type             = "Linux"
  sku_name            = "P1v3"
}

resource "azuresim_linux_web_app" "test" {
  name                = "app-linux"
  resource_group_name = "rg"
  location            = "eastus"
  service_plan_id     = azuresim_service_plan.plan.id
  https_only          = true

  app_settings = {
    NODE_ENV = "production"
  }

  site_config {
    always_on        = true
    linux_fx_version = "NODE|18-lts"
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_linux_web_app.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Web/sites/app-linux")),
					resource.TestCheckResourceAttr("azuresim_linux_web_app.test", "default_hostname",
						"app-linux.azurewebsites.net"),
					resource.TestCheckResourceAttrSet("azuresim_linux_web_app.test", "outbound_ip_addresses"),
				),
			},
		},
	})
}

func TestAccWindowsWebApp_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_windows_web_app" "test" {
  name                = "app-windows"
  resource_group_name = "rg"
  location            = "eastus"
  service_plan_id     = "/plans/p"

  site_config {
    always_on          = false
    windows_fx_version = "DOTNETCORE|8.0"
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_windows_web_app.test", "default_hostname",
					"app-windows.azurewebsites.net"),
			},
		},
	})
}

func TestAccLinuxFunctionApp_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_linux_function_app" "test" {
  name                       = "func-linux"
  resource_group_name        = "rg"
  location                   = "eastus"
  service_plan_id            = "/plans/p"
  storage_account_name       = "stfunc01"
  storage_account_access_key = "` + testFixtureStorageKey + `"
  https_only                 = true
  functions_extension_version = "~4"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_linux_function_app.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Web/sites/func-linux")),
					resource.TestCheckResourceAttr("azuresim_linux_function_app.test", "default_hostname",
						"func-linux.azurewebsites.net"),
				),
			},
		},
	})
}

func TestAccWindowsFunctionApp_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_windows_function_app" "test" {
  name                          = "func-windows"
  resource_group_name           = "rg"
  location                      = "eastus"
  service_plan_id               = "/plans/p"
  storage_account_name          = "stfunc02"
  storage_uses_managed_identity = true
}`,
				Check: resource.TestCheckResourceAttr("azuresim_windows_function_app.test", "id",
					arm("/resourceGroups/rg/providers/Microsoft.Web/sites/func-windows")),
			},
		},
	})
}
