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

// --- Monitor Action Group ---

var _ resource.Resource = &MonitorActionGroupResource{}
var _ resource.ResourceWithConfigure = &MonitorActionGroupResource{}

type MonitorActionGroupResource struct{ subscriptionID string }

type MonitorActionGroupModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	ShortName         types.String `tfsdk:"short_name"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	EmailReceiver     types.List   `tfsdk:"email_receiver"`
	SMSReceiver       types.List   `tfsdk:"sms_receiver"`
	WebhookReceiver   types.List   `tfsdk:"webhook_receiver"`
	Tags              types.Map    `tfsdk:"tags"`
}

func NewMonitorActionGroupResource() resource.Resource { return &MonitorActionGroupResource{} }

func (r *MonitorActionGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_action_group"
}

func (r *MonitorActionGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Monitor Action Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Action Group ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Action Group name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"short_name":         schema.StringAttribute{Required: true, Description: "Short name (max 12 chars)."},
			"enabled":            schema.BoolAttribute{Optional: true, Description: "Whether the action group is enabled."},
			"tags":               schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"email_receiver": schema.ListNestedBlock{
				Description: "Email receiver.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                    schema.StringAttribute{Required: true, Description: "Receiver name."},
						"email_address":           schema.StringAttribute{Required: true, Description: "Email address."},
						"use_common_alert_schema": schema.BoolAttribute{Optional: true, Description: "Use the common alert schema."},
					},
				},
			},
			"sms_receiver": schema.ListNestedBlock{
				Description: "SMS receiver.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":         schema.StringAttribute{Required: true, Description: "Receiver name."},
						"country_code": schema.StringAttribute{Required: true, Description: "Country code."},
						"phone_number": schema.StringAttribute{Required: true, Description: "Phone number."},
					},
				},
			},
			"webhook_receiver": schema.ListNestedBlock{
				Description: "Webhook receiver.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                    schema.StringAttribute{Required: true, Description: "Receiver name."},
						"service_uri":             schema.StringAttribute{Required: true, Description: "Webhook URL."},
						"use_common_alert_schema": schema.BoolAttribute{Optional: true, Description: "Use the common alert schema."},
					},
				},
			},
		},
	}
}

func (r *MonitorActionGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MonitorActionGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorActionGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/actionGroups/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorActionGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorActionGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorActionGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorActionGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state MonitorActionGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorActionGroupResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Monitor Metric Alert ---

var _ resource.Resource = &MonitorMetricAlertResource{}
var _ resource.ResourceWithConfigure = &MonitorMetricAlertResource{}

type MonitorMetricAlertResource struct{ subscriptionID string }

type MonitorMetricAlertModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	Description       types.String `tfsdk:"description"`
	Scopes            types.List   `tfsdk:"scopes"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Severity          types.Int64  `tfsdk:"severity"`
	Frequency         types.String `tfsdk:"frequency"`
	WindowSize        types.String `tfsdk:"window_size"`
	AutoMitigate      types.Bool   `tfsdk:"auto_mitigate"`
	TargetResourceType    types.String `tfsdk:"target_resource_type"`
	TargetResourceLocation types.String `tfsdk:"target_resource_location"`
	Criteria          types.List   `tfsdk:"criteria"`
	Action            types.List   `tfsdk:"action"`
	Tags              types.Map    `tfsdk:"tags"`
}

func NewMonitorMetricAlertResource() resource.Resource { return &MonitorMetricAlertResource{} }

func (r *MonitorMetricAlertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_metric_alert"
}

func (r *MonitorMetricAlertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Monitor metric alert.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Alert ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Alert name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"description":       schema.StringAttribute{Optional: true, Description: "Description."},
			"scopes":            schema.ListAttribute{Required: true, ElementType: types.StringType, Description: "Resource IDs to monitor."},
			"enabled":           schema.BoolAttribute{Optional: true, Description: "Whether the alert is enabled."},
			"severity":          schema.Int64Attribute{Optional: true, Description: "Severity (0-4)."},
			"frequency":         schema.StringAttribute{Optional: true, Description: "Evaluation frequency (ISO 8601, e.g. `PT1M`)."},
			"window_size":       schema.StringAttribute{Optional: true, Description: "Evaluation window (ISO 8601)."},
			"auto_mitigate":     schema.BoolAttribute{Optional: true, Description: "Auto-mitigate when condition clears."},
			"target_resource_type":     schema.StringAttribute{Optional: true, Description: "Target resource type (multi-resource alerts)."},
			"target_resource_location": schema.StringAttribute{Optional: true, Description: "Target resource region (multi-resource alerts)."},
			"tags":              schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"criteria": schema.ListNestedBlock{
				Description: "Static threshold criterion.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"metric_namespace":      schema.StringAttribute{Required: true, Description: "Metric namespace."},
						"metric_name":           schema.StringAttribute{Required: true, Description: "Metric name."},
						"aggregation":           schema.StringAttribute{Required: true, Description: "`Average`, `Count`, `Minimum`, `Maximum`, or `Total`."},
						"operator":              schema.StringAttribute{Required: true, Description: "`Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual`."},
						"threshold":             schema.Float64Attribute{Required: true, Description: "Threshold value."},
						"skip_metric_validation": schema.BoolAttribute{Optional: true, Description: "Skip Azure validation of the metric name."},
					},
				},
			},
			"action": schema.ListNestedBlock{
				Description: "Action linkage.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action_group_id":    schema.StringAttribute{Required: true, Description: "Action Group ID."},
						"webhook_properties": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Custom webhook properties."},
					},
				},
			},
		},
	}
}

func (r *MonitorMetricAlertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MonitorMetricAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorMetricAlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/metricAlerts/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorMetricAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorMetricAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorMetricAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorMetricAlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state MonitorMetricAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorMetricAlertResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
