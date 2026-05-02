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

// --- Container App Environment ---

var _ resource.Resource = &ContainerAppEnvironmentResource{}
var _ resource.ResourceWithConfigure = &ContainerAppEnvironmentResource{}

type ContainerAppEnvironmentResource struct{ subscriptionID string }

type ContainerAppEnvironmentModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	ResourceGroupName               types.String `tfsdk:"resource_group_name"`
	Location                        types.String `tfsdk:"location"`
	LogAnalyticsWorkspaceID         types.String `tfsdk:"log_analytics_workspace_id"`
	InfrastructureSubnetID          types.String `tfsdk:"infrastructure_subnet_id"`
	InternalLoadBalancerEnabled     types.Bool   `tfsdk:"internal_load_balancer_enabled"`
	ZoneRedundancyEnabled           types.Bool   `tfsdk:"zone_redundancy_enabled"`
	WorkloadProfileType             types.String `tfsdk:"workload_profile_type"`
	DefaultDomain                   types.String `tfsdk:"default_domain"`
	StaticIPAddress                 types.String `tfsdk:"static_ip_address"`
	Tags                            types.Map    `tfsdk:"tags"`
}

func NewContainerAppEnvironmentResource() resource.Resource { return &ContainerAppEnvironmentResource{} }

func (r *ContainerAppEnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_app_environment"
}

func (r *ContainerAppEnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Container App Environment.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Environment ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Environment name."},
			"resource_group_name":             schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":                        schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"log_analytics_workspace_id":      schema.StringAttribute{Optional: true, Description: "Linked LAW ID."},
			"infrastructure_subnet_id":        schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Infrastructure subnet ID."},
			"internal_load_balancer_enabled":  schema.BoolAttribute{Optional: true, Description: "Use an internal LB."},
			"zone_redundancy_enabled":         schema.BoolAttribute{Optional: true, Description: "Enable zone redundancy."},
			"workload_profile_type":           schema.StringAttribute{Optional: true, Description: "Workload profile (`Consumption`, `D4`, etc.)."},
			"default_domain": schema.StringAttribute{
				Computed: true, Description: "Simulated default domain.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"static_ip_address": schema.StringAttribute{
				Computed: true, Description: "Simulated static IP.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *ContainerAppEnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ContainerAppEnvironmentResource) applyComputed(plan *ContainerAppEnvironmentModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	loc := plan.Location.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s",
		r.subscriptionID, rg, name,
	))
	plan.DefaultDomain = types.StringValue(fmt.Sprintf("%s.%s.azurecontainerapps.io",
		strings.ToLower(name), strings.ToLower(loc)))
	plan.StaticIPAddress = types.StringValue("203.0.113.50")
}

func (r *ContainerAppEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerAppEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerAppEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerAppEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContainerAppEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerAppEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerAppEnvironmentResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Container App ---

var _ resource.Resource = &ContainerAppResource{}
var _ resource.ResourceWithConfigure = &ContainerAppResource{}

type ContainerAppResource struct{ subscriptionID string }

type ContainerAppModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	ContainerAppEnvironmentID types.String `tfsdk:"container_app_environment_id"`
	RevisionMode             types.String `tfsdk:"revision_mode"`
	WorkloadProfileName      types.String `tfsdk:"workload_profile_name"`
	Template                 types.List   `tfsdk:"template"`
	Ingress                  types.List   `tfsdk:"ingress"`
	LatestRevisionFQDN       types.String `tfsdk:"latest_revision_fqdn"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewContainerAppResource() resource.Resource { return &ContainerAppResource{} }

func (r *ContainerAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_app"
}

func (r *ContainerAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Container App.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Container App ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "App name."},
			"resource_group_name":          schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"container_app_environment_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent environment ID."},
			"revision_mode":                schema.StringAttribute{Required: true, Description: "`Single` or `Multiple`."},
			"workload_profile_name":        schema.StringAttribute{Optional: true, Description: "Workload profile name."},
			"latest_revision_fqdn": schema.StringAttribute{
				Computed: true, Description: "Simulated FQDN of the latest revision.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"template": schema.ListNestedBlock{
				Description: "Container template.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"min_replicas":            schema.Int64Attribute{Optional: true, Description: "Min replicas."},
						"max_replicas":            schema.Int64Attribute{Optional: true, Description: "Max replicas."},
						"revision_suffix":         schema.StringAttribute{Optional: true, Description: "Revision suffix."},
						"termination_grace_period_seconds": schema.Int64Attribute{Optional: true, Description: "Pod termination grace period."},
					},
					Blocks: map[string]schema.Block{
						"container": schema.ListNestedBlock{
							Description: "Container spec.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":   schema.StringAttribute{Required: true, Description: "Container name."},
									"image":  schema.StringAttribute{Required: true, Description: "Image (e.g. `nginx:latest`)."},
									"cpu":    schema.Float64Attribute{Required: true, Description: "vCPU (e.g. `0.25`)."},
									"memory": schema.StringAttribute{Required: true, Description: "Memory (e.g. `0.5Gi`)."},
									"args":    schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Container args."},
									"command": schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Container entrypoint."},
								},
							},
						},
					},
				},
			},
			"ingress": schema.ListNestedBlock{
				Description: "Ingress.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"external_enabled":      schema.BoolAttribute{Optional: true, Description: "Expose externally."},
						"target_port":           schema.Int64Attribute{Required: true, Description: "Target container port."},
						"transport":             schema.StringAttribute{Optional: true, Description: "`auto`, `http`, `http2`, or `tcp`."},
						"allow_insecure_connections": schema.BoolAttribute{Optional: true, Description: "Allow HTTP."},
					},
					Blocks: map[string]schema.Block{
						"traffic_weight": schema.ListNestedBlock{
							Description: "Traffic weight per revision.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"percentage":     schema.Int64Attribute{Required: true, Description: "Weight 0-100."},
									"latest_revision": schema.BoolAttribute{Optional: true, Description: "Apply to latest revision."},
									"revision_suffix": schema.StringAttribute{Optional: true, Description: "Specific revision suffix."},
									"label":          schema.StringAttribute{Optional: true, Description: "Traffic label."},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ContainerAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ContainerAppResource) applyComputed(plan *ContainerAppModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	plan.LatestRevisionFQDN = types.StringValue(fmt.Sprintf("%s.azurecontainerapps.io", plan.Name.ValueString()))
}

func (r *ContainerAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContainerAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerAppResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
