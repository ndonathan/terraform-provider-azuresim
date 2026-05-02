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

var _ resource.Resource = &RedisCacheResource{}
var _ resource.ResourceWithConfigure = &RedisCacheResource{}

type RedisCacheResource struct {
	subscriptionID string
}

type RedisCacheModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	ResourceGroupName          types.String `tfsdk:"resource_group_name"`
	Location                   types.String `tfsdk:"location"`
	Capacity                   types.Int64  `tfsdk:"capacity"`
	Family                     types.String `tfsdk:"family"`
	SKUName                    types.String `tfsdk:"sku_name"`
	NonSSLPortEnabled          types.Bool   `tfsdk:"non_ssl_port_enabled"`
	MinimumTLSVersion          types.String `tfsdk:"minimum_tls_version"`
	RedisVersion               types.String `tfsdk:"redis_version"`
	SubnetID                   types.String `tfsdk:"subnet_id"`
	PrivateStaticIPAddress     types.String `tfsdk:"private_static_ip_address"`
	Zones                      types.List   `tfsdk:"zones"`
	PublicNetworkAccessEnabled types.Bool   `tfsdk:"public_network_access_enabled"`
	Hostname                   types.String `tfsdk:"hostname"`
	Port                       types.Int64  `tfsdk:"port"`
	SSLPort                    types.Int64  `tfsdk:"ssl_port"`
	PrimaryAccessKey           types.String `tfsdk:"primary_access_key"`
	SecondaryAccessKey         types.String `tfsdk:"secondary_access_key"`
	PrimaryConnectionString    types.String `tfsdk:"primary_connection_string"`
	SecondaryConnectionString  types.String `tfsdk:"secondary_connection_string"`
	Tags                       types.Map    `tfsdk:"tags"`
}

func NewRedisCacheResource() resource.Resource { return &RedisCacheResource{} }

func (r *RedisCacheResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redis_cache"
}

func (r *RedisCacheResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Cache for Redis instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Redis Cache ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the cache (globally unique).",
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
			"capacity":                       schema.Int64Attribute{Required: true, Description: "Cache size (`Basic`/`Standard`: 0-6, `Premium`: 1-5)."},
			"family":                         schema.StringAttribute{Required: true, Description: "`C` (Basic/Standard) or `P` (Premium)."},
			"sku_name":                       schema.StringAttribute{Required: true, Description: "`Basic`, `Standard`, or `Premium`."},
			"non_ssl_port_enabled":           schema.BoolAttribute{Optional: true, Description: "Enable the non-SSL port (6379)."},
			"minimum_tls_version":            schema.StringAttribute{Optional: true, Description: "`1.0`, `1.1`, or `1.2`."},
			"redis_version":                  schema.StringAttribute{Optional: true, Description: "Redis version (`4` or `6`)."},
			"subnet_id":                      schema.StringAttribute{Optional: true, Description: "Subnet ID (Premium only)."},
			"private_static_ip_address":      schema.StringAttribute{Optional: true, Description: "Static IP address inside the subnet."},
			"zones":                          schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
			"public_network_access_enabled":  schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"hostname": schema.StringAttribute{
				Computed: true, Description: "Simulated hostname (`<name>.redis.cache.windows.net`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"port": schema.Int64Attribute{Computed: true, Description: "Non-SSL port (6379)."},
			"ssl_port": schema.Int64Attribute{Computed: true, Description: "SSL port (6380)."},
			"primary_access_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary access key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secondary_access_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary access key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"primary_connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary connection string.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secondary_connection_string": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary connection string.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *RedisCacheResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func mkRedisKey(seed string) string {
	sum := sha1.Sum([]byte(seed))
	padded := make([]byte, 32)
	for i := range padded {
		padded[i] = sum[i%len(sum)]
	}
	return base64.StdEncoding.EncodeToString(padded)
}

func (r *RedisCacheResource) applyComputed(plan *RedisCacheModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/Redis/%s",
		r.subscriptionID, rg, name,
	))
	plan.Hostname = types.StringValue(fmt.Sprintf("%s.redis.cache.windows.net", name))
	plan.Port = types.Int64Value(6379)
	plan.SSLPort = types.Int64Value(6380)

	primary := mkRedisKey("redis-primary/" + rg + "/" + name)
	secondary := mkRedisKey("redis-secondary/" + rg + "/" + name)
	plan.PrimaryAccessKey = types.StringValue(primary)
	plan.SecondaryAccessKey = types.StringValue(secondary)
	plan.PrimaryConnectionString = types.StringValue(fmt.Sprintf(
		"%s.redis.cache.windows.net:6380,password=%s,ssl=True,abortConnect=False", name, primary,
	))
	plan.SecondaryConnectionString = types.StringValue(fmt.Sprintf(
		"%s.redis.cache.windows.net:6380,password=%s,ssl=True,abortConnect=False", name, secondary,
	))
}

func (r *RedisCacheResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RedisCacheModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RedisCacheResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RedisCacheModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RedisCacheResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RedisCacheModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state RedisCacheModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.Hostname = state.Hostname
	plan.Port = state.Port
	plan.SSLPort = state.SSLPort
	plan.PrimaryAccessKey = state.PrimaryAccessKey
	plan.SecondaryAccessKey = state.SecondaryAccessKey
	plan.PrimaryConnectionString = state.PrimaryConnectionString
	plan.SecondaryConnectionString = state.SecondaryConnectionString
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RedisCacheResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
