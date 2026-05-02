package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMSSQLServerAndDatabase(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_mssql_server" "srv" {
  name                         = "sql-server-001"
  resource_group_name          = "rg"
  location                     = "eastus"
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "` + testFixturePassword + `"
  minimum_tls_version          = "1.2"
}

resource "azuresim_mssql_database" "db" {
  name      = "appdb"
  server_id = azuresim_mssql_server.srv.id
  sku_name  = "S0"
  collation = "SQL_Latin1_General_CP1_CI_AS"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_mssql_server.srv", "fully_qualified_domain_name",
						"sql-server-001.database.windows.net"),
					resource.TestCheckResourceAttr("azuresim_mssql_database.db", "id",
						arm("/resourceGroups/rg/providers/Microsoft.Sql/servers/sql-server-001/databases/appdb")),
				),
			},
		},
	})
}

func TestAccCosmosDBAccount_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_cosmosdb_account" "test" {
  name                = "cosmos-001"
  resource_group_name = "rg"
  location            = "eastus"
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = "eastus"
    failover_priority = 0
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_cosmosdb_account.test", "endpoint",
						"https://cosmos-001.documents.azure.com:443/"),
					resource.TestCheckResourceAttrSet("azuresim_cosmosdb_account.test", "primary_key"),
					resource.TestCheckResourceAttrSet("azuresim_cosmosdb_account.test", "primary_readonly_key"),
				),
			},
		},
	})
}

func TestAccRedisCache_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_redis_cache" "test" {
  name                = "redis-001"
  resource_group_name = "rg"
  location            = "eastus"
  capacity            = 1
  family              = "C"
  sku_name            = "Standard"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_redis_cache.test", "hostname",
						"redis-001.redis.cache.windows.net"),
					resource.TestCheckResourceAttr("azuresim_redis_cache.test", "ssl_port", "6380"),
					resource.TestCheckResourceAttr("azuresim_redis_cache.test", "port", "6379"),
					resource.TestMatchResourceAttr("azuresim_redis_cache.test", "primary_connection_string",
						regexp.MustCompile(`redis-001\.redis\.cache\.windows\.net:6380,password=`)),
				),
			},
		},
	})
}

func TestAccPostgresFlexible_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_postgresql_flexible_server" "test" {
  name                = "pg-001"
  resource_group_name = "rg"
  location            = "eastus"
  version             = "16"
  sku_name            = "GP_Standard_D2s_v3"
  storage_mb          = 32768
  administrator_login    = "pgadmin"
  administrator_password = "` + testFixturePassword + `"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_postgresql_flexible_server.test", "fqdn",
					"pg-001.postgres.database.azure.com"),
			},
		},
	})
}

func TestAccMySQLFlexible_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_mysql_flexible_server" "test" {
  name                = "mysql-001"
  resource_group_name = "rg"
  location            = "eastus"
  version             = "8.0.21"
  sku_name            = "B_Standard_B1ms"

  storage {
    size_gb           = 32
    auto_grow_enabled = true
  }
}`,
				Check: resource.TestCheckResourceAttr("azuresim_mysql_flexible_server.test", "fqdn",
					"mysql-001.mysql.database.azure.com"),
			},
		},
	})
}
