package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageAccount_basicAndComputed(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_storage_account" "test" {
  name                     = "stexample0001"
  resource_group_name      = "rg"
  location                 = "eastus"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_storage_account.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/stexample0001")),
					resource.TestCheckResourceAttr("azuresim_storage_account.test", "primary_blob_endpoint",
						"https://stexample0001.blob.core.windows.net/"),
					resource.TestCheckResourceAttrSet("azuresim_storage_account.test", "primary_access_key"),
				),
			},
		},
	})
}

func TestAccStorageContainer_basicAndDefaultAccessType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_storage_container" "test" {
  name                 = "data"
  storage_account_name = "stexample0001"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_storage_container.test", "id",
						"https://stexample0001.blob.core.windows.net/data"),
					resource.TestCheckResourceAttr("azuresim_storage_container.test", "container_access_type", "private"),
					resource.TestCheckResourceAttr("azuresim_storage_container.test", "has_immutability_policy", "false"),
					resource.TestCheckResourceAttr("azuresim_storage_container.test", "has_legal_hold", "false"),
				),
			},
		},
	})
}

func TestAccStorageBlob_url(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_storage_blob" "test" {
  name                   = "config.json"
  storage_account_name   = "stexample0001"
  storage_container_name = "data"
  type                   = "BlockBlob"
  source_content         = "{}"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_storage_blob.test", "id",
						"https://stexample0001.blob.core.windows.net/data/config.json"),
					resource.TestCheckResourceAttr("azuresim_storage_blob.test", "url",
						"https://stexample0001.blob.core.windows.net/data/config.json"),
				),
			},
		},
	})
}
