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

var _ resource.Resource = &RouteTableResource{}
var _ resource.ResourceWithConfigure = &RouteTableResource{}

type RouteTableResource struct {
	subscriptionID string
}

type RouteTableModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	BGPRoutePropagationEnabled  types.Bool   `tfsdk:"bgp_route_propagation_enabled"`
	Route                       types.List   `tfsdk:"route"`
	Tags                        types.Map    `tfsdk:"tags"`
}

type RouteModel struct {
	Name             types.String `tfsdk:"name"`
	AddressPrefix    types.String `tfsdk:"address_prefix"`
	NextHopType      types.String `tfsdk:"next_hop_type"`
	NextHopInIPAddress types.String `tfsdk:"next_hop_in_ip_address"`
}

func NewRouteTableResource() resource.Resource {
	return &RouteTableResource{}
}

func (r *RouteTableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_table"
}

func (r *RouteTableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Route Table.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Route Table ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the Route Table.",
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
			"bgp_route_propagation_enabled": schema.BoolAttribute{
				Optional: true, Description: "Whether BGP route propagation is enabled. Defaults to `true`.",
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
		Blocks: map[string]schema.Block{
			"route": schema.ListNestedBlock{
				Description: "Inline route entry.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":           schema.StringAttribute{Required: true, Description: "Route name."},
						"address_prefix": schema.StringAttribute{Required: true, Description: "Destination CIDR."},
						"next_hop_type":  schema.StringAttribute{Required: true, Description: "`VirtualNetworkGateway`, `VnetLocal`, `Internet`, `VirtualAppliance`, or `None`."},
						"next_hop_in_ip_address": schema.StringAttribute{
							Optional: true, Description: "Required when `next_hop_type` is `VirtualAppliance`.",
						},
					},
				},
			},
		},
	}
}

func (r *RouteTableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *RouteTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RouteTableModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RouteTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RouteTableModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state RouteTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RouteTableResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Subnet Route Table Association ---

var _ resource.Resource = &SubnetRouteTableAssociationResource{}
var _ resource.ResourceWithConfigure = &SubnetRouteTableAssociationResource{}

type SubnetRouteTableAssociationResource struct {
	subscriptionID string
}

type SubnetRouteTableAssociationModel struct {
	ID           types.String `tfsdk:"id"`
	SubnetID     types.String `tfsdk:"subnet_id"`
	RouteTableID types.String `tfsdk:"route_table_id"`
}

func NewSubnetRouteTableAssociationResource() resource.Resource {
	return &SubnetRouteTableAssociationResource{}
}

func (r *SubnetRouteTableAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet_route_table_association"
}

func (r *SubnetRouteTableAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associates a Route Table with a Subnet. Mirrors the AzureRM convention of using the Subnet ID as the resource ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Subnet ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"subnet_id": schema.StringAttribute{
				Required: true, Description: "The Subnet ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"route_table_id": schema.StringAttribute{Required: true, Description: "The Route Table ID."},
		},
	}
}

func (r *SubnetRouteTableAssociationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *SubnetRouteTableAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubnetRouteTableAssociationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = plan.SubnetID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetRouteTableAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubnetRouteTableAssociationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SubnetRouteTableAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SubnetRouteTableAssociationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = plan.SubnetID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetRouteTableAssociationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
