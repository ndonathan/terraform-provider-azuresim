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

var _ resource.Resource = &PrivateDNSZoneResource{}
var _ resource.ResourceWithConfigure = &PrivateDNSZoneResource{}

type PrivateDNSZoneResource struct {
	subscriptionID string
}

type PrivateDNSZoneModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	ResourceGroupName               types.String `tfsdk:"resource_group_name"`
	NumberOfRecordSets              types.Int64  `tfsdk:"number_of_record_sets"`
	MaxNumberOfRecordSets           types.Int64  `tfsdk:"max_number_of_record_sets"`
	NumberOfVirtualNetworkLinks     types.Int64  `tfsdk:"number_of_virtual_network_links"`
	MaxNumberOfVirtualNetworkLinks  types.Int64  `tfsdk:"max_number_of_virtual_network_links"`
	Tags                            types.Map    `tfsdk:"tags"`
}

func NewPrivateDNSZoneResource() resource.Resource {
	return &PrivateDNSZoneResource{}
}

func (r *PrivateDNSZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_dns_zone"
}

func (r *PrivateDNSZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Private DNS Zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Private DNS Zone ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Zone name (e.g. `privatelink.blob.core.windows.net`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"number_of_record_sets": schema.Int64Attribute{
				Computed: true, Description: "Always `0` in this simulator.",
			},
			"max_number_of_record_sets": schema.Int64Attribute{
				Computed: true, Description: "Static (`25000`).",
			},
			"number_of_virtual_network_links": schema.Int64Attribute{
				Computed: true, Description: "Always `0` in this simulator.",
			},
			"max_number_of_virtual_network_links": schema.Int64Attribute{
				Computed: true, Description: "Static (`1000`).",
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *PrivateDNSZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *PrivateDNSZoneResource) applyComputed(plan *PrivateDNSZoneModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	plan.NumberOfRecordSets = types.Int64Value(0)
	plan.MaxNumberOfRecordSets = types.Int64Value(25000)
	plan.NumberOfVirtualNetworkLinks = types.Int64Value(0)
	plan.MaxNumberOfVirtualNetworkLinks = types.Int64Value(1000)
}

func (r *PrivateDNSZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PrivateDNSZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateDNSZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PrivateDNSZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateDNSZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PrivateDNSZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateDNSZoneResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Private DNS Zone VNet Link ---

var _ resource.Resource = &PrivateDNSZoneVNetLinkResource{}
var _ resource.ResourceWithConfigure = &PrivateDNSZoneVNetLinkResource{}

type PrivateDNSZoneVNetLinkResource struct {
	subscriptionID string
}

type PrivateDNSZoneVNetLinkModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	PrivateDNSZoneName  types.String `tfsdk:"private_dns_zone_name"`
	VirtualNetworkID    types.String `tfsdk:"virtual_network_id"`
	RegistrationEnabled types.Bool   `tfsdk:"registration_enabled"`
	Tags                types.Map    `tfsdk:"tags"`
}

func NewPrivateDNSZoneVNetLinkResource() resource.Resource {
	return &PrivateDNSZoneVNetLinkResource{}
}

func (r *PrivateDNSZoneVNetLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_dns_zone_virtual_network_link"
}

func (r *PrivateDNSZoneVNetLinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Links a Virtual Network to a Private DNS Zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Link ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Link name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group of the parent Private DNS Zone.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"private_dns_zone_name": schema.StringAttribute{
				Required: true, Description: "Parent Private DNS Zone name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"virtual_network_id": schema.StringAttribute{
				Required: true, Description: "Virtual Network ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"registration_enabled": schema.BoolAttribute{
				Optional: true, Description: "Enable auto-registration of VM hostnames.",
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *PrivateDNSZoneVNetLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *PrivateDNSZoneVNetLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PrivateDNSZoneVNetLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/virtualNetworkLinks/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(),
		plan.PrivateDNSZoneName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateDNSZoneVNetLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PrivateDNSZoneVNetLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PrivateDNSZoneVNetLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PrivateDNSZoneVNetLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state PrivateDNSZoneVNetLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PrivateDNSZoneVNetLinkResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
