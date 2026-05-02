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

var _ resource.Resource = &ApplicationInsightsResource{}
var _ resource.ResourceWithConfigure = &ApplicationInsightsResource{}

type ApplicationInsightsResource struct {
	subscriptionID string
}

type ApplicationInsightsModel struct {
	ID                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	ResourceGroupName                   types.String `tfsdk:"resource_group_name"`
	Location                            types.String `tfsdk:"location"`
	ApplicationType                     types.String `tfsdk:"application_type"`
	WorkspaceID                         types.String `tfsdk:"workspace_id"`
	RetentionInDays                     types.Int64  `tfsdk:"retention_in_days"`
	SamplingPercentage                  types.Float64 `tfsdk:"sampling_percentage"`
	DailyDataCapInGB                    types.Float64 `tfsdk:"daily_data_cap_in_gb"`
	DailyDataCapNotificationsDisabled   types.Bool   `tfsdk:"daily_data_cap_notifications_disabled"`
	DisableIPMasking                    types.Bool   `tfsdk:"disable_ip_masking"`
	LocalAuthenticationDisabled         types.Bool   `tfsdk:"local_authentication_disabled"`
	InternetIngestionEnabled            types.Bool   `tfsdk:"internet_ingestion_enabled"`
	InternetQueryEnabled                types.Bool   `tfsdk:"internet_query_enabled"`
	ForceCustomerStorageForProfiler     types.Bool   `tfsdk:"force_customer_storage_for_profiler"`
	AppID                               types.String `tfsdk:"app_id"`
	InstrumentationKey                  types.String `tfsdk:"instrumentation_key"`
	ConnectionString                    types.String `tfsdk:"connection_string"`
	Tags                                types.Map    `tfsdk:"tags"`
}

func NewApplicationInsightsResource() resource.Resource {
	return &ApplicationInsightsResource{}
}

func (r *ApplicationInsightsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_insights"
}

func (r *ApplicationInsightsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Application Insights component.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Application Insights ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the component.",
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
			"application_type": schema.StringAttribute{
				Required: true, Description: "`web`, `other`, `java`, `MobileCenter`, `Node.JS`, `phone`, or `store`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"workspace_id": schema.StringAttribute{
				Optional: true, Description: "Linked Log Analytics Workspace ID (for workspace-based instances).",
			},
			"retention_in_days":                     schema.Int64Attribute{Optional: true, Description: "Data retention (30, 60, 90, 120, 180, 270, 365, 550, 730)."},
			"sampling_percentage":                   schema.Float64Attribute{Optional: true, Description: "Sampling percentage (0-100)."},
			"daily_data_cap_in_gb":                  schema.Float64Attribute{Optional: true, Description: "Daily data cap in GB."},
			"daily_data_cap_notifications_disabled": schema.BoolAttribute{Optional: true, Description: "Disable email notifications when daily cap is hit."},
			"disable_ip_masking":                    schema.BoolAttribute{Optional: true, Description: "Disable IP masking."},
			"local_authentication_disabled":         schema.BoolAttribute{Optional: true, Description: "Disable non-AAD authentication."},
			"internet_ingestion_enabled":            schema.BoolAttribute{Optional: true, Description: "Allow public internet ingestion."},
			"internet_query_enabled":                schema.BoolAttribute{Optional: true, Description: "Allow public internet query."},
			"force_customer_storage_for_profiler":   schema.BoolAttribute{Optional: true, Description: "Force customer storage for Profiler."},
			"app_id": schema.StringAttribute{
				Computed: true, Description: "Simulated AppID (UUID).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"instrumentation_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated instrumentation key (UUID).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated connection string.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *ApplicationInsightsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ApplicationInsightsResource) applyComputed(plan *ApplicationInsightsModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s",
		r.subscriptionID, rg, name,
	))
	appID := simulatedUUID("appInsights/appId/" + rg + "/" + name)
	iKey := simulatedUUID("appInsights/iKey/" + rg + "/" + name)
	plan.AppID = types.StringValue(appID)
	plan.InstrumentationKey = types.StringValue(iKey)
	plan.ConnectionString = types.StringValue(fmt.Sprintf(
		"InstrumentationKey=%s;IngestionEndpoint=https://%s.in.applicationinsights.azure.com/;LiveEndpoint=https://%s.livediagnostics.monitor.azure.com/",
		iKey, plan.Location.ValueString(), plan.Location.ValueString(),
	))
}

func (r *ApplicationInsightsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationInsightsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationInsightsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationInsightsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationInsightsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationInsightsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ApplicationInsightsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.AppID = state.AppID
	plan.InstrumentationKey = state.InstrumentationKey
	plan.ConnectionString = state.ConnectionString
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationInsightsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
