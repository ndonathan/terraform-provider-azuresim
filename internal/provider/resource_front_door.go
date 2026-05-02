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

var _ resource.Resource = &FrontDoorProfileResource{}
var _ resource.ResourceWithConfigure = &FrontDoorProfileResource{}

type FrontDoorProfileResource struct{ subscriptionID string }

type FrontDoorProfileModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	SKUName                  types.String `tfsdk:"sku_name"`
	ResponseTimeoutSeconds   types.Int64  `tfsdk:"response_timeout_seconds"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewFrontDoorProfileResource() resource.Resource { return &FrontDoorProfileResource{} }

func (r *FrontDoorProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_frontdoor_profile"
}

func (r *FrontDoorProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Front Door (Standard/Premium) profile.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Profile ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Profile name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"sku_name":               schema.StringAttribute{Required: true, Description: "`Standard_AzureFrontDoor` or `Premium_AzureFrontDoor`."},
			"response_timeout_seconds": schema.Int64Attribute{Optional: true, Description: "Origin response timeout (16-240)."},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *FrontDoorProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *FrontDoorProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FrontDoorProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FrontDoorProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FrontDoorProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FrontDoorProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FrontDoorProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state FrontDoorProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FrontDoorProfileResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
