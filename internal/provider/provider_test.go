package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	testSubscriptionID = "12345678-1234-1234-1234-123456789012"
	testTenantID       = "87654321-4321-4321-4321-210987654321"
)

// providerConfig is prepended to every test config to wire the provider with
// a deterministic subscription/tenant ID. Resources reference these in their
// generated IDs so they can be asserted against without env-dependent values.
const providerConfig = `
provider "azuresim" {
  subscription_id = "` + testSubscriptionID + `"
  tenant_id       = "` + testTenantID + `"
}
`

// testAccProtoV6ProviderFactories wires the in-process provider into the
// terraform-plugin-testing framework. Tests run with `resource.UnitTest` so
// they execute without any external setup or `TF_ACC`.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"azuresim": providerserver.NewProtocol6WithError(New("test")()),
}
