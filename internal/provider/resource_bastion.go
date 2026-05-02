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

var _ resource.Resource = &BastionHostResource{}
var _ resource.ResourceWithConfigure = &BastionHostResource{}

type BastionHostResource struct{ subscriptionID string }

type BastionHostModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	Location                 types.String `tfsdk:"location"`
	SKU                      types.String `tfsdk:"sku"`
	ScaleUnits               types.Int64  `tfsdk:"scale_units"`
	CopyPasteEnabled         types.Bool   `tfsdk:"copy_paste_enabled"`
	FileCopyEnabled          types.Bool   `tfsdk:"file_copy_enabled"`
	IPConnectEnabled         types.Bool   `tfsdk:"ip_connect_enabled"`
	ShareableLinkEnabled     types.Bool   `tfsdk:"shareable_link_enabled"`
	TunnelingEnabled         types.Bool   `tfsdk:"tunneling_enabled"`
	IPConfiguration          types.List   `tfsdk:"ip_configuration"`
	DNSName                  types.String `tfsdk:"dns_name"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewBastionHostResource() resource.Resource { return &BastionHostResource{} }

func (r *BastionHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bastion_host"
}

func (r *BastionHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Bastion Host.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Bastion ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Bastion name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":              schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"sku":                   schema.StringAttribute{Optional: true, Description: "`Basic`, `Standard`, `Developer`, or `Premium`."},
			"scale_units":           schema.Int64Attribute{Optional: true, Description: "Scale units (Standard only, 2-50)."},
			"copy_paste_enabled":    schema.BoolAttribute{Optional: true, Description: "Allow copy/paste."},
			"file_copy_enabled":     schema.BoolAttribute{Optional: true, Description: "Allow file copy (Standard only)."},
			"ip_connect_enabled":    schema.BoolAttribute{Optional: true, Description: "Allow IP-based connect."},
			"shareable_link_enabled": schema.BoolAttribute{Optional: true, Description: "Allow shareable links."},
			"tunneling_enabled":     schema.BoolAttribute{Optional: true, Description: "Allow native client tunneling."},
			"dns_name": schema.StringAttribute{
				Computed: true, Description: "Simulated DNS name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"ip_configuration": schema.ListNestedBlock{
				Description: "IP configuration. Subnet must be named `AzureBastionSubnet`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                 schema.StringAttribute{Required: true, Description: "Configuration name."},
						"subnet_id":            schema.StringAttribute{Required: true, Description: "`AzureBastionSubnet` ID."},
						"public_ip_address_id": schema.StringAttribute{Required: true, Description: "Public IP ID."},
					},
				},
			},
		},
	}
}

func (r *BastionHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *BastionHostResource) applyComputed(plan *BastionHostModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/bastionHosts/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.DNSName = types.StringValue(fmt.Sprintf("bst-%s.bastion.azure.com", simulatedUUID("bastion/"+name)[:8]))
}

func (r *BastionHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BastionHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BastionHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BastionHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BastionHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BastionHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state BastionHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.DNSName = state.DNSName
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BastionHostResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
