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

var _ resource.Resource = &EventHubNamespaceResource{}
var _ resource.ResourceWithConfigure = &EventHubNamespaceResource{}

type EventHubNamespaceResource struct {
	subscriptionID string
}

type EventHubNamespaceModel struct {
	ID                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	ResourceGroupName                   types.String `tfsdk:"resource_group_name"`
	Location                            types.String `tfsdk:"location"`
	SKU                                 types.String `tfsdk:"sku"`
	Capacity                            types.Int64  `tfsdk:"capacity"`
	AutoInflateEnabled                  types.Bool   `tfsdk:"auto_inflate_enabled"`
	MaximumThroughputUnits              types.Int64  `tfsdk:"maximum_throughput_units"`
	ZoneRedundant                       types.Bool   `tfsdk:"zone_redundant"`
	MinimumTLSVersion                   types.String `tfsdk:"minimum_tls_version"`
	PublicNetworkAccessEnabled          types.Bool   `tfsdk:"public_network_access_enabled"`
	LocalAuthEnabled                    types.Bool   `tfsdk:"local_auth_enabled"`
	DefaultPrimaryConnectionString      types.String `tfsdk:"default_primary_connection_string"`
	DefaultSecondaryConnectionString    types.String `tfsdk:"default_secondary_connection_string"`
	DefaultPrimaryKey                   types.String `tfsdk:"default_primary_key"`
	DefaultSecondaryKey                 types.String `tfsdk:"default_secondary_key"`
	Tags                                types.Map    `tfsdk:"tags"`
}

func NewEventHubNamespaceResource() resource.Resource { return &EventHubNamespaceResource{} }

func (r *EventHubNamespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventhub_namespace"
}

func (r *EventHubNamespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Event Hubs Namespace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Namespace ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Namespace name (globally unique).",
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
			"sku":                            schema.StringAttribute{Required: true, Description: "`Basic`, `Standard`, or `Premium`."},
			"capacity":                       schema.Int64Attribute{Optional: true, Description: "Throughput units."},
			"auto_inflate_enabled":           schema.BoolAttribute{Optional: true, Description: "Enable auto-inflate (Standard only)."},
			"maximum_throughput_units":       schema.Int64Attribute{Optional: true, Description: "Auto-inflate ceiling."},
			"zone_redundant":                 schema.BoolAttribute{Optional: true, Description: "Enable zone redundancy."},
			"minimum_tls_version":            schema.StringAttribute{Optional: true, Description: "`1.0`, `1.1`, or `1.2`."},
			"public_network_access_enabled":  schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"local_auth_enabled":             schema.BoolAttribute{Optional: true, Description: "Enable SAS-key auth."},
			"default_primary_connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary connection string.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_secondary_connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary connection string.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_primary_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_secondary_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *EventHubNamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *EventHubNamespaceResource) applyComputed(plan *EventHubNamespaceModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s",
		r.subscriptionID, rg, name,
	))
	primary := mkSBKey("eh-primary/" + rg + "/" + name)
	secondary := mkSBKey("eh-secondary/" + rg + "/" + name)
	plan.DefaultPrimaryKey = types.StringValue(primary)
	plan.DefaultSecondaryKey = types.StringValue(secondary)
	plan.DefaultPrimaryConnectionString = types.StringValue(fmt.Sprintf(
		"Endpoint=sb://%s.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=%s",
		name, primary,
	))
	plan.DefaultSecondaryConnectionString = types.StringValue(fmt.Sprintf(
		"Endpoint=sb://%s.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=%s",
		name, secondary,
	))
}

func (r *EventHubNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EventHubNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EventHubNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EventHubNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EventHubNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state EventHubNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.DefaultPrimaryKey = state.DefaultPrimaryKey
	plan.DefaultSecondaryKey = state.DefaultSecondaryKey
	plan.DefaultPrimaryConnectionString = state.DefaultPrimaryConnectionString
	plan.DefaultSecondaryConnectionString = state.DefaultSecondaryConnectionString
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubNamespaceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
