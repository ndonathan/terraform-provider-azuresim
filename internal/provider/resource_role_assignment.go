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

var _ resource.Resource = &RoleAssignmentResource{}
var _ resource.ResourceWithConfigure = &RoleAssignmentResource{}

type RoleAssignmentResource struct {
	subscriptionID string
}

type RoleAssignmentModel struct {
	ID                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Scope                          types.String `tfsdk:"scope"`
	RoleDefinitionID               types.String `tfsdk:"role_definition_id"`
	RoleDefinitionName             types.String `tfsdk:"role_definition_name"`
	PrincipalID                    types.String `tfsdk:"principal_id"`
	PrincipalType                  types.String `tfsdk:"principal_type"`
	Description                    types.String `tfsdk:"description"`
	Condition                      types.String `tfsdk:"condition"`
	ConditionVersion               types.String `tfsdk:"condition_version"`
	SkipServicePrincipalAADCheck   types.Bool   `tfsdk:"skip_service_principal_aad_check"`
	DelegatedManagedIdentityResource types.String `tfsdk:"delegated_managed_identity_resource_id"`
}

func NewRoleAssignmentResource() resource.Resource {
	return &RoleAssignmentResource{}
}

func (r *RoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_assignment"
}

func (r *RoleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure RBAC Role Assignment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Role Assignment ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional: true, Computed: true,
				Description: "Assignment name (UUID). Generated when omitted.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"scope": schema.StringAttribute{
				Required: true, Description: "Resource scope (e.g. a subscription, RG, or resource ID).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"role_definition_id": schema.StringAttribute{
				Optional: true, Description: "Role definition ID. Mutually exclusive with `role_definition_name`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"role_definition_name": schema.StringAttribute{
				Optional: true, Description: "Role definition name (e.g. `Contributor`). Mutually exclusive with `role_definition_id`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"principal_id": schema.StringAttribute{
				Required: true, Description: "Object ID of the principal.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"principal_type": schema.StringAttribute{
				Optional: true, Computed: true,
				Description: "Principal type (`User`, `Group`, `ServicePrincipal`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description":                            schema.StringAttribute{Optional: true, Description: "Description for the assignment."},
			"condition":                              schema.StringAttribute{Optional: true, Description: "ABAC condition expression."},
			"condition_version":                      schema.StringAttribute{Optional: true, Description: "Condition version (e.g. `2.0`)."},
			"skip_service_principal_aad_check":       schema.BoolAttribute{Optional: true, Description: "Skip AAD propagation check (no-op in sim)."},
			"delegated_managed_identity_resource_id": schema.StringAttribute{Optional: true, Description: "Delegated managed identity resource ID."},
		},
	}
}

func (r *RoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *RoleAssignmentResource) applyComputed(plan *RoleAssignmentModel) {
	if plan.Name.IsNull() || plan.Name.IsUnknown() || plan.Name.ValueString() == "" {
		seed := plan.Scope.ValueString() + "/" + plan.PrincipalID.ValueString() + "/" +
			plan.RoleDefinitionID.ValueString() + plan.RoleDefinitionName.ValueString()
		plan.Name = types.StringValue(simulatedUUID("roleAssignment/" + seed))
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"%s/providers/Microsoft.Authorization/roleAssignments/%s",
		plan.Scope.ValueString(), plan.Name.ValueString(),
	))
	if plan.PrincipalType.IsNull() || plan.PrincipalType.IsUnknown() || plan.PrincipalType.ValueString() == "" {
		plan.PrincipalType = types.StringValue("ServicePrincipal")
	}
}

func (r *RoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state RoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.Name = state.Name
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RoleAssignmentResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
