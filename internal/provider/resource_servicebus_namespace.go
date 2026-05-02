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

var _ resource.Resource = &ServiceBusNamespaceResource{}
var _ resource.ResourceWithConfigure = &ServiceBusNamespaceResource{}

type ServiceBusNamespaceResource struct {
	subscriptionID string
}

type ServiceBusNamespaceModel struct {
	ID                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	ResourceGroupName                   types.String `tfsdk:"resource_group_name"`
	Location                            types.String `tfsdk:"location"`
	SKU                                 types.String `tfsdk:"sku"`
	Capacity                            types.Int64  `tfsdk:"capacity"`
	PremiumMessagingPartitions          types.Int64  `tfsdk:"premium_messaging_partitions"`
	MinimumTLSVersion                   types.String `tfsdk:"minimum_tls_version"`
	PublicNetworkAccessEnabled          types.Bool   `tfsdk:"public_network_access_enabled"`
	LocalAuthEnabled                    types.Bool   `tfsdk:"local_auth_enabled"`
	Endpoint                            types.String `tfsdk:"endpoint"`
	DefaultPrimaryConnectionString      types.String `tfsdk:"default_primary_connection_string"`
	DefaultSecondaryConnectionString    types.String `tfsdk:"default_secondary_connection_string"`
	DefaultPrimaryKey                   types.String `tfsdk:"default_primary_key"`
	DefaultSecondaryKey                 types.String `tfsdk:"default_secondary_key"`
	Tags                                types.Map    `tfsdk:"tags"`
}

func NewServiceBusNamespaceResource() resource.Resource { return &ServiceBusNamespaceResource{} }

func (r *ServiceBusNamespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicebus_namespace"
}

func (r *ServiceBusNamespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Service Bus Namespace.",
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
			"sku":                          schema.StringAttribute{Required: true, Description: "`Basic`, `Standard`, or `Premium`."},
			"capacity":                     schema.Int64Attribute{Optional: true, Description: "Premium messaging units (1, 2, 4, 8, 16)."},
			"premium_messaging_partitions": schema.Int64Attribute{Optional: true, Description: "Premium partitions (1, 2, 4)."},
			"minimum_tls_version":          schema.StringAttribute{Optional: true, Description: "`1.0`, `1.1`, or `1.2`."},
			"public_network_access_enabled": schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"local_auth_enabled":           schema.BoolAttribute{Optional: true, Description: "Enable SAS-key auth."},
			"endpoint": schema.StringAttribute{
				Computed: true, Description: "Simulated SBus endpoint.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_primary_connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary connection string for `RootManageSharedAccessKey`.",
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

func (r *ServiceBusNamespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func mkSBKey(seed string) string {
	sum := sha1.Sum([]byte(seed))
	padded := make([]byte, 32)
	for i := range padded {
		padded[i] = sum[i%len(sum)]
	}
	return base64.StdEncoding.EncodeToString(padded)
}

func (r *ServiceBusNamespaceResource) applyComputed(plan *ServiceBusNamespaceModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s",
		r.subscriptionID, rg, name,
	))
	plan.Endpoint = types.StringValue(fmt.Sprintf("sb://%s.servicebus.windows.net/", name))

	primary := mkSBKey("sb-primary/" + rg + "/" + name)
	secondary := mkSBKey("sb-secondary/" + rg + "/" + name)
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

func (r *ServiceBusNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceBusNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceBusNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceBusNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceBusNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ServiceBusNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.Endpoint = state.Endpoint
	plan.DefaultPrimaryKey = state.DefaultPrimaryKey
	plan.DefaultSecondaryKey = state.DefaultSecondaryKey
	plan.DefaultPrimaryConnectionString = state.DefaultPrimaryConnectionString
	plan.DefaultSecondaryConnectionString = state.DefaultSecondaryConnectionString
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusNamespaceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
