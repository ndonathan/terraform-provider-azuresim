package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FirewallResource{}
var _ resource.ResourceWithConfigure = &FirewallResource{}

type FirewallResource struct{ subscriptionID string }

type FirewallModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	Location            types.String `tfsdk:"location"`
	SKUName             types.String `tfsdk:"sku_name"`
	SKUTier             types.String `tfsdk:"sku_tier"`
	FirewallPolicyID    types.String `tfsdk:"firewall_policy_id"`
	ThreatIntelMode     types.String `tfsdk:"threat_intel_mode"`
	DNSServers          types.List   `tfsdk:"dns_servers"`
	PrivateIPRanges     types.List   `tfsdk:"private_ip_ranges"`
	Zones               types.List   `tfsdk:"zones"`
	IPConfiguration     types.List   `tfsdk:"ip_configuration"`
	ManagementIPConfig  types.List   `tfsdk:"management_ip_configuration"`
	Tags                types.Map    `tfsdk:"tags"`
}

func NewFirewallResource() resource.Resource { return &FirewallResource{} }

func (r *FirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

func (r *FirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	ipConfig := schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"name":                 schema.StringAttribute{Required: true, Description: "Configuration name."},
			"subnet_id":            schema.StringAttribute{Optional: true, Description: "`AzureFirewallSubnet` ID."},
			"public_ip_address_id": schema.StringAttribute{Optional: true, Description: "Public IP ID."},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Firewall.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Firewall ID."},
			"name":                schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Firewall name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":            schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"sku_name":            schema.StringAttribute{Required: true, Description: "`AZFW_VNet` or `AZFW_Hub`."},
			"sku_tier":            schema.StringAttribute{Required: true, Description: "`Standard`, `Premium`, or `Basic`."},
			"firewall_policy_id":  schema.StringAttribute{Optional: true, Description: "Firewall policy ID."},
			"threat_intel_mode":   schema.StringAttribute{Optional: true, Description: "`Off`, `Alert`, or `Deny`."},
			"dns_servers":         schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Custom DNS servers."},
			"private_ip_ranges":   schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Private IP ranges (SNAT exemptions)."},
			"zones":               schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
			"tags":                schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"ip_configuration":            schema.ListNestedBlock{Description: "Frontend IP configurations.", NestedObject: ipConfig},
			"management_ip_configuration": schema.ListNestedBlock{Description: "Management IP configuration (Forced Tunneling).", NestedObject: ipConfig},
		},
	}
}

func (r *FirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *FirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state FirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
