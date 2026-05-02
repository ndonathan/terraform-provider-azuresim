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

var _ resource.Resource = &VirtualNetworkPeeringResource{}
var _ resource.ResourceWithConfigure = &VirtualNetworkPeeringResource{}

type VirtualNetworkPeeringResource struct {
	subscriptionID string
}

type VirtualNetworkPeeringModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	ResourceGroupName         types.String `tfsdk:"resource_group_name"`
	VirtualNetworkName        types.String `tfsdk:"virtual_network_name"`
	RemoteVirtualNetworkID    types.String `tfsdk:"remote_virtual_network_id"`
	AllowVirtualNetworkAccess types.Bool   `tfsdk:"allow_virtual_network_access"`
	AllowForwardedTraffic     types.Bool   `tfsdk:"allow_forwarded_traffic"`
	AllowGatewayTransit       types.Bool   `tfsdk:"allow_gateway_transit"`
	UseRemoteGateways         types.Bool   `tfsdk:"use_remote_gateways"`
}

func NewVirtualNetworkPeeringResource() resource.Resource {
	return &VirtualNetworkPeeringResource{}
}

func (r *VirtualNetworkPeeringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_network_peering"
}

func (r *VirtualNetworkPeeringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates a Virtual Network Peering. Create one resource per direction.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Peering ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Peering name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group of the local VNet.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"virtual_network_name": schema.StringAttribute{
				Required: true, Description: "Local VNet name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"remote_virtual_network_id": schema.StringAttribute{
				Required: true, Description: "Remote VNet ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"allow_virtual_network_access": schema.BoolAttribute{Optional: true, Description: "Allow access to remote VNet."},
			"allow_forwarded_traffic":      schema.BoolAttribute{Optional: true, Description: "Allow forwarded traffic from the remote VNet."},
			"allow_gateway_transit":        schema.BoolAttribute{Optional: true, Description: "Allow gateway transit."},
			"use_remote_gateways":          schema.BoolAttribute{Optional: true, Description: "Use the remote VNet's gateway."},
		},
	}
}

func (r *VirtualNetworkPeeringResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *VirtualNetworkPeeringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VirtualNetworkPeeringModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/virtualNetworkPeerings/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(),
		plan.VirtualNetworkName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualNetworkPeeringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualNetworkPeeringModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualNetworkPeeringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VirtualNetworkPeeringModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state VirtualNetworkPeeringModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualNetworkPeeringResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
