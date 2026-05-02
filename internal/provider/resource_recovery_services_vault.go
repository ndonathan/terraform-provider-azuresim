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

var _ resource.Resource = &RecoveryServicesVaultResource{}
var _ resource.ResourceWithConfigure = &RecoveryServicesVaultResource{}

type RecoveryServicesVaultResource struct {
	subscriptionID string
}

type RecoveryServicesVaultModel struct {
	ID                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	ResourceGroupName             types.String `tfsdk:"resource_group_name"`
	Location                      types.String `tfsdk:"location"`
	SKU                           types.String `tfsdk:"sku"`
	StorageModeType               types.String `tfsdk:"storage_mode_type"`
	SoftDeleteEnabled             types.Bool   `tfsdk:"soft_delete_enabled"`
	PublicNetworkAccessEnabled    types.Bool   `tfsdk:"public_network_access_enabled"`
	CrossRegionRestoreEnabled     types.Bool   `tfsdk:"cross_region_restore_enabled"`
	ImmutabilityEnabled           types.String `tfsdk:"immutability"`
	Tags                          types.Map    `tfsdk:"tags"`
}

func NewRecoveryServicesVaultResource() resource.Resource { return &RecoveryServicesVaultResource{} }

func (r *RecoveryServicesVaultResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recovery_services_vault"
}

func (r *RecoveryServicesVaultResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Recovery Services Vault.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Vault ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Vault name.",
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
			"sku":                          schema.StringAttribute{Required: true, Description: "`Standard` or `RS0`."},
			"storage_mode_type":            schema.StringAttribute{Optional: true, Description: "`GeoRedundant`, `LocallyRedundant`, or `ZoneRedundant`."},
			"soft_delete_enabled":          schema.BoolAttribute{Optional: true, Description: "Enable soft delete."},
			"public_network_access_enabled": schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"cross_region_restore_enabled": schema.BoolAttribute{Optional: true, Description: "Enable cross-region restore (requires `GeoRedundant`)."},
			"immutability":                 schema.StringAttribute{Optional: true, Description: "`Disabled`, `Unlocked`, or `Locked`."},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *RecoveryServicesVaultResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *RecoveryServicesVaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RecoveryServicesVaultModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RecoveryServicesVaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RecoveryServicesVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RecoveryServicesVaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RecoveryServicesVaultModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state RecoveryServicesVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RecoveryServicesVaultResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
