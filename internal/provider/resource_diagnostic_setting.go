package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &MonitorDiagnosticSettingResource{}
var _ resource.ResourceWithConfigure = &MonitorDiagnosticSettingResource{}

type MonitorDiagnosticSettingResource struct{ subscriptionID string }

type MonitorDiagnosticSettingModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	TargetResourceID            types.String `tfsdk:"target_resource_id"`
	LogAnalyticsWorkspaceID     types.String `tfsdk:"log_analytics_workspace_id"`
	StorageAccountID            types.String `tfsdk:"storage_account_id"`
	EventHubAuthorizationRuleID types.String `tfsdk:"eventhub_authorization_rule_id"`
	EventHubName                types.String `tfsdk:"eventhub_name"`
	LogAnalyticsDestinationType types.String `tfsdk:"log_analytics_destination_type"`
	EnabledLog                  types.List   `tfsdk:"enabled_log"`
	Metric                      types.List   `tfsdk:"metric"`
}

func NewMonitorDiagnosticSettingResource() resource.Resource { return &MonitorDiagnosticSettingResource{} }

func (r *MonitorDiagnosticSettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_diagnostic_setting"
}

func (r *MonitorDiagnosticSettingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Monitor diagnostic setting on a target resource.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Diagnostic Setting ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Diagnostic Setting name."},
			"target_resource_id":          schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Target resource ID."},
			"log_analytics_workspace_id":  schema.StringAttribute{Optional: true, Description: "Log Analytics Workspace destination."},
			"storage_account_id":          schema.StringAttribute{Optional: true, Description: "Storage Account destination."},
			"eventhub_authorization_rule_id": schema.StringAttribute{Optional: true, Description: "Event Hub authorization rule ID destination."},
			"eventhub_name":                  schema.StringAttribute{Optional: true, Description: "Event Hub name destination."},
			"log_analytics_destination_type": schema.StringAttribute{Optional: true, Description: "`Dedicated` or `AzureDiagnostics`."},
		},
		Blocks: map[string]schema.Block{
			"enabled_log": schema.ListNestedBlock{
				Description: "Enabled log category.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"category":       schema.StringAttribute{Optional: true, Description: "Log category."},
						"category_group": schema.StringAttribute{Optional: true, Description: "Log category group (`allLogs`, `audit`, etc.)."},
					},
				},
			},
			"metric": schema.ListNestedBlock{
				Description: "Metric category.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{Required: true, Description: "Metric category (`AllMetrics` is most common)."},
						"enabled":  schema.BoolAttribute{Optional: true, Description: "Whether this metric is enabled."},
					},
				},
			},
		},
	}
}

func (r *MonitorDiagnosticSettingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MonitorDiagnosticSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorDiagnosticSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s|%s",
		strings.TrimRight(plan.TargetResourceID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorDiagnosticSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorDiagnosticSettingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorDiagnosticSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorDiagnosticSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s|%s",
		strings.TrimRight(plan.TargetResourceID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorDiagnosticSettingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
