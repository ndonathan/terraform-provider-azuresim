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

var _ resource.Resource = &PrivateEndpointResource{}
var _ resource.ResourceWithConfigure = &PrivateEndpointResource{}

type PrivateEndpointResource struct{ subscriptionID string }

type PrivateEndpointModel struct {
	ID                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	ResourceGroupName             types.String `tfsdk:"resource_group_name"`
	Location                      types.String `tfsdk:"location"`
	SubnetID                      types.String `tfsdk:"subnet_id"`
	CustomNetworkInterfaceName    types.String `tfsdk:"custom_network_interface_name"`
	PrivateServiceConnection      types.List   `tfsdk:"private_service_connection"`
	PrivateDNSZoneGroup           types.List   `tfsdk:"private_dns_zone_group"`
	Tags                          types.Map    `tfsdk:"tags"`
}

func NewPrivateEndpointResource() resource.Resource { return &PrivateEndpointResource{} }

func (r *PrivateEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoint"
}

func (r *PrivateEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Private Endpoint.",
		Attributes: map[string]schema.Attribute{
			"id":                            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Private Endpoint ID."},
			"name":                          schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Endpoint name."},
			"resource_group_name":           schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":                      schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"subnet_id":                     schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Subnet ID."},
			"custom_network_interface_name": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Custom NIC name."},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"private_service_connection": schema.ListNestedBlock{
				Description: "Service connection.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                              schema.StringAttribute{Required: true, Description: "Connection name."},
						"private_connection_resource_id":    schema.StringAttribute{Optional: true, Description: "Target resource ID."},
						"private_connection_resource_alias": schema.StringAttribute{Optional: true, Description: "Target resource alias."},
						"is_manual_connection":              schema.BoolAttribute{Required: true, Description: "Whether the connection requires manual approval."},
						"subresource_names":                 schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Subresource names (e.g. `[\"blob\"]`)."},
						"request_message":                   schema.StringAttribute{Optional: true, Description: "Approval request message (manual only)."},
					},
				},
			},
			"private_dns_zone_group": schema.ListNestedBlock{
				Description: "Optional DNS zone group binding.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                 schema.StringAttribute{Required: true, Description: "Group name."},
						"private_dns_zone_ids": schema.ListAttribute{Required: true, ElementType: types.StringType, Description: "Private DNS zone IDs."},
					},
				},
			},
		},
	}
}

func (r *PrivateEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *PrivateEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PrivateEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PrivateEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PrivateEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state PrivateEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateEndpointResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
