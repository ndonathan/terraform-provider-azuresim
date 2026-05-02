package provider

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &UserAssignedIdentityResource{}
var _ resource.ResourceWithConfigure = &UserAssignedIdentityResource{}

type UserAssignedIdentityResource struct {
	subscriptionID string
}

type UserAssignedIdentityModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	Location          types.String `tfsdk:"location"`
	PrincipalID       types.String `tfsdk:"principal_id"`
	ClientID          types.String `tfsdk:"client_id"`
	TenantID          types.String `tfsdk:"tenant_id"`
	Tags              types.Map    `tfsdk:"tags"`
}

func NewUserAssignedIdentityResource() resource.Resource {
	return &UserAssignedIdentityResource{}
}

func (r *UserAssignedIdentityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_assigned_identity"
}

func (r *UserAssignedIdentityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure User-Assigned Managed Identity.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The simulated Managed Identity ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the User-Assigned Identity.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Name of the Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"principal_id": schema.StringAttribute{
				Computed: true, Description: "Simulated service principal object ID (UUID).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"client_id": schema.StringAttribute{
				Computed: true, Description: "Simulated client (application) ID (UUID).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tenant_id": schema.StringAttribute{
				Computed: true, Description: "Tenant ID under which the identity is created.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *UserAssignedIdentityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

// simulatedUUID derives a deterministic v4-shaped UUID from a seed.
// Output format matches the canonical 8-4-4-4-12 hex layout.
func simulatedUUID(seed string) string {
	sum := sha1.Sum([]byte(seed))
	// Force version (4) and variant (10xx) bits to match v4 UUID shape.
	sum[6] = (sum[6] & 0x0f) | 0x40
	sum[8] = (sum[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

func (r *UserAssignedIdentityResource) applyComputed(plan *UserAssignedIdentityModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s",
		r.subscriptionID, rg, name,
	))
	plan.PrincipalID = types.StringValue(simulatedUUID("principal/" + rg + "/" + name))
	plan.ClientID = types.StringValue(simulatedUUID("client/" + rg + "/" + name))
	plan.TenantID = types.StringValue("00000000-0000-0000-0000-000000000000")
}

func (r *UserAssignedIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserAssignedIdentityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserAssignedIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserAssignedIdentityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserAssignedIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserAssignedIdentityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state UserAssignedIdentityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.PrincipalID = state.PrincipalID
	plan.ClientID = state.ClientID
	plan.TenantID = state.TenantID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserAssignedIdentityResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
