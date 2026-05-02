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

var _ resource.Resource = &ServicePlanResource{}
var _ resource.ResourceWithConfigure = &ServicePlanResource{}

type ServicePlanResource struct {
	subscriptionID string
}

type ServicePlanModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	Location                 types.String `tfsdk:"location"`
	OSType                   types.String `tfsdk:"os_type"`
	SKUName                  types.String `tfsdk:"sku_name"`
	WorkerCount              types.Int64  `tfsdk:"worker_count"`
	MaximumElasticWorkerCount types.Int64 `tfsdk:"maximum_elastic_worker_count"`
	PerSiteScalingEnabled    types.Bool   `tfsdk:"per_site_scaling_enabled"`
	ZoneBalancingEnabled     types.Bool   `tfsdk:"zone_balancing_enabled"`
	AppServiceEnvironmentID  types.String `tfsdk:"app_service_environment_id"`
	Reserved                 types.Bool   `tfsdk:"reserved"`
	Kind                     types.String `tfsdk:"kind"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewServicePlanResource() resource.Resource {
	return &ServicePlanResource{}
}

func (r *ServicePlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_plan"
}

func (r *ServicePlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure App Service Plan (`Microsoft.Web/serverfarms`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Service Plan ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the Service Plan.",
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
			"os_type": schema.StringAttribute{
				Required: true, Description: "`Linux`, `Windows`, or `WindowsContainer`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"sku_name":                     schema.StringAttribute{Required: true, Description: "SKU (e.g. `B1`, `S1`, `P1v3`, `Y1`, `EP1`)."},
			"worker_count":                 schema.Int64Attribute{Optional: true, Description: "Number of workers."},
			"maximum_elastic_worker_count": schema.Int64Attribute{Optional: true, Description: "Max elastic worker count (Premium plans)."},
			"per_site_scaling_enabled":     schema.BoolAttribute{Optional: true, Description: "Enable per-site scaling."},
			"zone_balancing_enabled":       schema.BoolAttribute{Optional: true, Description: "Enable zone balancing."},
			"app_service_environment_id":   schema.StringAttribute{Optional: true, Description: "ASE ID for ASEv3 plans."},
			"reserved": schema.BoolAttribute{
				Computed: true, Description: "Whether the plan is Linux (`reserved=true`).",
				PlanModifiers: []planmodifier.Bool{},
			},
			"kind": schema.StringAttribute{
				Computed: true, Description: "Plan kind (e.g. `linux`, `app`, `functionapp`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *ServicePlanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ServicePlanResource) applyComputed(plan *ServicePlanModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverfarms/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	switch plan.OSType.ValueString() {
	case "Linux":
		plan.Reserved = types.BoolValue(true)
		plan.Kind = types.StringValue("linux")
	case "WindowsContainer":
		plan.Reserved = types.BoolValue(false)
		plan.Kind = types.StringValue("windows,container")
	default:
		plan.Reserved = types.BoolValue(false)
		plan.Kind = types.StringValue("app")
	}
}

func (r *ServicePlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServicePlanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServicePlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServicePlanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServicePlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServicePlanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ServicePlanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServicePlanResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
