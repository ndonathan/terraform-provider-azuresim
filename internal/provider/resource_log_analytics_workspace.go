package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &LogAnalyticsWorkspaceResource{}
var _ resource.ResourceWithConfigure = &LogAnalyticsWorkspaceResource{}

type LogAnalyticsWorkspaceResource struct {
	subscriptionID string
}

type LogAnalyticsWorkspaceModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	SKU                         types.String `tfsdk:"sku"`
	RetentionInDays             types.Int64  `tfsdk:"retention_in_days"`
	DailyQuotaGB                types.Float64 `tfsdk:"daily_quota_gb"`
	InternetIngestionEnabled    types.Bool   `tfsdk:"internet_ingestion_enabled"`
	InternetQueryEnabled        types.Bool   `tfsdk:"internet_query_enabled"`
	ReservationCapacityGBPerDay types.Int64  `tfsdk:"reservation_capacity_in_gb_per_day"`
	LocalAuthenticationDisabled types.Bool   `tfsdk:"local_authentication_disabled"`
	CMKForQueryForced           types.Bool   `tfsdk:"cmk_for_query_forced"`
	WorkspaceID                 types.String `tfsdk:"workspace_id"`
	PrimarySharedKey            types.String `tfsdk:"primary_shared_key"`
	SecondarySharedKey          types.String `tfsdk:"secondary_shared_key"`
	Tags                        types.Map    `tfsdk:"tags"`
}

func NewLogAnalyticsWorkspaceResource() resource.Resource {
	return &LogAnalyticsWorkspaceResource{}
}

func (r *LogAnalyticsWorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_analytics_workspace"
}

func (r *LogAnalyticsWorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Log Analytics Workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The simulated Workspace ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the workspace.",
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
			"sku": schema.StringAttribute{
				Optional: true, Computed: true,
				Description: "SKU name. Defaults to `PerGB2018`. Other options include `Free`, `Standalone`, `PerNode`, `CapacityReservation`.",
			},
			"retention_in_days": schema.Int64Attribute{
				Optional: true, Computed: true,
				Description: "Retention in days (30-730). Defaults to 30.",
			},
			"daily_quota_gb":                     schema.Float64Attribute{Optional: true, Description: "Daily ingestion cap in GB. Use `-1` for unlimited."},
			"internet_ingestion_enabled":         schema.BoolAttribute{Optional: true, Description: "Whether public internet ingestion is allowed."},
			"internet_query_enabled":             schema.BoolAttribute{Optional: true, Description: "Whether public internet query is allowed."},
			"reservation_capacity_in_gb_per_day": schema.Int64Attribute{Optional: true, Description: "Reserved capacity (only with `CapacityReservation` SKU)."},
			"local_authentication_disabled":      schema.BoolAttribute{Optional: true, Description: "Disable non-AAD authentication."},
			"cmk_for_query_forced":               schema.BoolAttribute{Optional: true, Description: "Force CMK for query."},
			"workspace_id": schema.StringAttribute{
				Computed: true, Description: "Simulated Workspace (Customer) ID — UUID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"primary_shared_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary shared key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secondary_shared_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary shared key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *LogAnalyticsWorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func simulatedSharedKey(seed string) string {
	sum := sha1.Sum([]byte(seed))
	// Pad to 64 bytes for a realistic-looking key, then base64-encode.
	padded := make([]byte, 64)
	for i := range padded {
		padded[i] = sum[i%len(sum)]
	}
	return base64.StdEncoding.EncodeToString(padded)
}

func (r *LogAnalyticsWorkspaceResource) applyComputed(plan *LogAnalyticsWorkspaceModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s",
		r.subscriptionID, rg, name,
	))
	plan.WorkspaceID = types.StringValue(simulatedUUID("workspace/" + rg + "/" + name))
	plan.PrimarySharedKey = types.StringValue(simulatedSharedKey("primary/" + rg + "/" + name))
	plan.SecondarySharedKey = types.StringValue(simulatedSharedKey("secondary/" + rg + "/" + name))

	if plan.SKU.IsNull() || plan.SKU.IsUnknown() {
		plan.SKU = types.StringValue("PerGB2018")
	}
	if plan.RetentionInDays.IsNull() || plan.RetentionInDays.IsUnknown() {
		plan.RetentionInDays = types.Int64Value(30)
	}
}

func (r *LogAnalyticsWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LogAnalyticsWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogAnalyticsWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LogAnalyticsWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LogAnalyticsWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LogAnalyticsWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state LogAnalyticsWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.WorkspaceID = state.WorkspaceID
	plan.PrimarySharedKey = state.PrimarySharedKey
	plan.SecondarySharedKey = state.SecondarySharedKey
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LogAnalyticsWorkspaceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
