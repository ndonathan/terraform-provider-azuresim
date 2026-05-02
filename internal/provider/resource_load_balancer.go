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

var _ resource.Resource = &LoadBalancerResource{}
var _ resource.ResourceWithConfigure = &LoadBalancerResource{}

type LoadBalancerResource struct {
	subscriptionID string
}

type LoadBalancerModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	Location                 types.String `tfsdk:"location"`
	SKU                      types.String `tfsdk:"sku"`
	SKUTier                  types.String `tfsdk:"sku_tier"`
	EdgeZone                 types.String `tfsdk:"edge_zone"`
	FrontendIPConfiguration  types.List   `tfsdk:"frontend_ip_configuration"`
	PrivateIPAddress         types.String `tfsdk:"private_ip_address"`
	PrivateIPAddresses       types.List   `tfsdk:"private_ip_addresses"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewLoadBalancerResource() resource.Resource { return &LoadBalancerResource{} }

func (r *LoadBalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lb"
}

func (r *LoadBalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Load Balancer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Load Balancer ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Load Balancer name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"sku":      schema.StringAttribute{Optional: true, Description: "`Basic`, `Standard`, or `Gateway`. Defaults to `Basic`."},
			"sku_tier": schema.StringAttribute{Optional: true, Description: "`Regional` or `Global`."},
			"edge_zone": schema.StringAttribute{
				Optional: true, Description: "Edge zone.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"private_ip_address": schema.StringAttribute{
				Computed: true, Description: "Primary private IP (mirrors first frontend with a private allocation).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"private_ip_addresses": schema.ListAttribute{
				Computed: true, ElementType: types.StringType, Description: "All private IPs.",
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"frontend_ip_configuration": schema.ListNestedBlock{
				Description: "Frontend IP configuration. At least one is required.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                          schema.StringAttribute{Required: true, Description: "Configuration name."},
						"public_ip_address_id":          schema.StringAttribute{Optional: true, Description: "Public IP ID (for public LBs)."},
						"subnet_id":                     schema.StringAttribute{Optional: true, Description: "Subnet ID (for internal LBs)."},
						"private_ip_address":            schema.StringAttribute{Optional: true, Description: "Static private IP (when allocation is `Static`)."},
						"private_ip_address_allocation": schema.StringAttribute{Optional: true, Description: "`Static` or `Dynamic`."},
						"private_ip_address_version":    schema.StringAttribute{Optional: true, Description: "`IPv4` or `IPv6`."},
						"zones":                         schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
					},
				},
			},
		},
	}
}

func (r *LoadBalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *LoadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LoadBalancerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	plan.PrivateIPAddress = types.StringValue("")
	plan.PrivateIPAddresses = types.ListValueMust(types.StringType, nil)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LoadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LoadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LoadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LoadBalancerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state LoadBalancerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.PrivateIPAddress = state.PrivateIPAddress
	plan.PrivateIPAddresses = state.PrivateIPAddresses
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LoadBalancerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
