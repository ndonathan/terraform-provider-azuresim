package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &AzureSimProvider{}

type AzureSimProvider struct {
	version string
}

type AzureSimProviderModel struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AzureSimProvider{
			version: version,
		}
	}
}

func (p *AzureSimProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azuresim"
	resp.Version = p.version
}

func (p *AzureSimProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A simulated Azure provider for testing Terraform configurations without provisioning real resources.",
		Attributes: map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				Optional:    true,
				Description: "The simulated Azure Subscription ID.",
			},
			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The simulated Azure Tenant ID.",
			},
		},
	}
}

func (p *AzureSimProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config AzureSimProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriptionID := "00000000-0000-0000-0000-000000000000"
	if !config.SubscriptionID.IsNull() && !config.SubscriptionID.IsUnknown() {
		subscriptionID = config.SubscriptionID.ValueString()
	}

	// Pass subscription ID to resources via provider data
	resp.DataSourceData = subscriptionID
	resp.ResourceData = subscriptionID
}

func (p *AzureSimProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceGroupResource,
		NewVirtualNetworkResource,
		NewSubnetResource,
		NewVirtualMachineResource,
		NewStorageAccountResource,
	}
}

func (p *AzureSimProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
