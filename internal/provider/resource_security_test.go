package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserAssignedIdentity_computedUUIDs(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_user_assigned_identity" "test" {
  name                = "uai-app"
  resource_group_name = "rg"
  location            = "eastus"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_user_assigned_identity.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uai-app")),
					resource.TestMatchResourceAttr("azuresim_user_assigned_identity.test", "principal_id", uuidV4Regex),
					resource.TestMatchResourceAttr("azuresim_user_assigned_identity.test", "client_id", uuidV4Regex),
				),
			},
		},
	})
}

func TestAccUserAssignedIdentity_principalIDStableAcrossUpdates(t *testing.T) {
	// Verifies that a tag-only update doesn't perturb the deterministic UUIDs.
	// Captures principal_id from step 1 and asserts it survives into step 2.
	var savedPrincipalID string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_user_assigned_identity" "test" {
  name                = "uai-stable"
  resource_group_name = "rg"
  location            = "eastus"
  tags                = { v = "1" }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_user_assigned_identity.test", "principal_id", uuidV4Regex),
					func(s *terraform.State) error {
						savedPrincipalID = s.RootModule().Resources["azuresim_user_assigned_identity.test"].Primary.Attributes["principal_id"]
						return nil
					},
				),
			},
			{
				Config: providerConfig + `
resource "azuresim_user_assigned_identity" "test" {
  name                = "uai-stable"
  resource_group_name = "rg"
  location            = "eastus"
  tags                = { v = "2" }
}`,
				Check: func(s *terraform.State) error {
					got := s.RootModule().Resources["azuresim_user_assigned_identity.test"].Primary.Attributes["principal_id"]
					if got != savedPrincipalID {
						return fmt.Errorf("principal_id changed across update: %q -> %q", savedPrincipalID, got)
					}
					return nil
				},
			},
		},
	})
}

func TestAccKeyVault_basicAndDefaults(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_key_vault" "test" {
  name                = "kv-test-001"
  resource_group_name = "rg"
  location            = "eastus"
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_key_vault.test", "id",
						arm("/resourceGroups/rg/providers/Microsoft.KeyVault/vaults/kv-test-001")),
					resource.TestCheckResourceAttr("azuresim_key_vault.test", "vault_uri", "https://kv-test-001.vault.azure.net/"),
					resource.TestCheckResourceAttr("azuresim_key_vault.test", "soft_delete_retention_days", "90"),
					resource.TestCheckResourceAttr("azuresim_key_vault.test", "public_network_access_enabled", "true"),
				),
			},
		},
	})
}

func TestAccKeyVaultSecret_versionedURIs(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_key_vault" "kv" {
  name                = "kv-secrets"
  resource_group_name = "rg"
  location            = "eastus"
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_key_vault_secret" "test" {
  name         = "db-password"
  key_vault_id = azuresim_key_vault.kv.id
  value        = "p@ssword"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_key_vault_secret.test", "version",
						regexp.MustCompile(`^[0-9a-f]{32}$`)),
					resource.TestMatchResourceAttr("azuresim_key_vault_secret.test", "id",
						regexp.MustCompile(`^https://kv-secrets\.vault\.azure\.net/secrets/db-password/[0-9a-f]{32}$`)),
					resource.TestCheckResourceAttr("azuresim_key_vault_secret.test", "versionless_id",
						"https://kv-secrets.vault.azure.net/secrets/db-password"),
				),
			},
		},
	})
}

func TestAccKeyVaultKey_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_key_vault" "kv" {
  name                = "kv-keys"
  resource_group_name = "rg"
  location            = "eastus"
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "azuresim_key_vault_key" "test" {
  name         = "data-key"
  key_vault_id = azuresim_key_vault.kv.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["encrypt", "decrypt"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_key_vault_key.test", "id",
						regexp.MustCompile(`^https://kv-keys\.vault\.azure\.net/keys/data-key/[0-9a-f]{32}$`)),
					resource.TestCheckResourceAttrSet("azuresim_key_vault_key.test", "public_key_pem"),
				),
			},
		},
	})
}

func TestAccRoleAssignment_generatedName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_role_assignment" "test" {
  scope                = "/subscriptions/abc/resourceGroups/rg"
  role_definition_name = "Contributor"
  principal_id         = "11111111-1111-1111-1111-111111111111"
  principal_type       = "ServicePrincipal"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("azuresim_role_assignment.test", "name", uuidV4Regex),
					resource.TestMatchResourceAttr("azuresim_role_assignment.test", "id",
						regexp.MustCompile(`^/subscriptions/abc/resourceGroups/rg/providers/Microsoft\.Authorization/roleAssignments/[0-9a-f]{8}-`)),
				),
			},
		},
	})
}

func TestAccRoleAssignment_explicitNamePersists(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_role_assignment" "test" {
  name                 = "00000000-0000-0000-0000-000000000001"
  scope                = "/subscriptions/abc"
  role_definition_name = "Reader"
  principal_id         = "22222222-2222-2222-2222-222222222222"
}`,
				Check: resource.TestCheckResourceAttr("azuresim_role_assignment.test", "id",
					"/subscriptions/abc/providers/Microsoft.Authorization/roleAssignments/00000000-0000-0000-0000-000000000001"),
			},
		},
	})
}
